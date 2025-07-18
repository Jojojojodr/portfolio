package views

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Jojojojodr/portfolio/frontend/auth"
	"github.com/Jojojojodr/portfolio/frontend/components"
	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *gin.Context) {
    name := c.PostForm("name")
    password := c.PostForm("password")

    if name == "" || password == "" {
		c.Writer.WriteHeader(400)
		components.LoginResponse("", "Username and password required").Render(c.Request.Context(), c.Writer)
        return
    }

    user, err := models.GetUserByName(name)
    if err != nil || user == nil {
		c.Writer.WriteHeader(401)
		components.LoginResponse("", "Invalid username or password").Render(c.Request.Context(), c.Writer)
        return
    }

    if !internal.CheckPasswordHash(password, user.Password) {
		c.Writer.WriteHeader(401)
		components.LoginResponse("", "Invalid username or password").Render(c.Request.Context(), c.Writer)
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
    })

    secretToken := internal.Env("SECRET_TOKEN")
    tokenString, err := token.SignedString([]byte(secretToken))
    if err != nil {
		c.Writer.WriteHeader(500)
		components.LoginResponse("", "Could not create token").Render(c.Request.Context(), c.Writer)
        return
    }

    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)

	c.Header("HX-Location", "/")
	c.Writer.WriteHeader(200)
}

func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "/", "", false, true) // Clear the cookie

	c.Redirect(http.StatusSeeOther, "/")
}

func ProfileHandler(c *gin.Context) {
    currentUser, exists := c.Get("user")
    if !exists {
        c.Redirect(http.StatusFound, "/login")
        return
    }

    user := currentUser.(*models.User)
    userID := c.Param("id")
    isAdmin := c.GetBool("isAdmin")
    
    var targetUser *models.User
    var isOwnProfile bool

    if userID == "" {
        // Viewing own profile
        targetUser = user
        isOwnProfile = true
    } else {
        // Viewing another user's profile (admin only)
        if !isAdmin {
            c.AbortWithStatus(http.StatusForbidden)
            return
        }
        
        // Fetch the target user from database
        id, err := strconv.Atoi(userID)
        if err != nil {
            c.AbortWithStatus(http.StatusBadRequest)
            return
        }
        
        targetUser, err = models.GetUserByID(uint(id))
        if err != nil || targetUser == nil {
            c.AbortWithStatus(http.StatusNotFound)
            return
        }
        
        isOwnProfile = false
    }

    errors := make(map[string]string)
    
    // Set proper content type for HTML response
    c.Header("Content-Type", "text/html; charset=utf-8")
    
    component := auth.ProfilePage(c, targetUser, isOwnProfile, errors)
    err := component.Render(c.Request.Context(), c.Writer)
    if err != nil {
        c.String(http.StatusInternalServerError, "Error rendering profile page: %v", err)
        return
    }
}

func UpdateProfileHandler(c *gin.Context) {
    // Get current user from session
    currentUser, exists := c.Get("user")
    if !exists {
        c.Redirect(http.StatusFound, "/login")
        return
    }

    user := currentUser.(*models.User)
    isAdmin := c.GetBool("isAdmin")
    
    // Get form data
    userIDStr := c.PostForm("user_id")
    name := strings.TrimSpace(c.PostForm("name"))
    email := strings.TrimSpace(c.PostForm("email"))
    currentPassword := c.PostForm("current_password")
    newPassword := c.PostForm("new_password")
    confirmPassword := c.PostForm("confirm_password")
    isAdminStr := c.PostForm("is_admin") // This will be "on" if checked, empty if not
    
    // Determine target user
    var targetUser *models.User
    var isOwnProfile bool
    
    if userIDStr == "" || userIDStr == strconv.Itoa(int(user.ID)) {
        // Updating own profile
        targetUser = user
        isOwnProfile = true
    } else {
        // Admin updating another user's profile
        if !isAdmin {
            c.AbortWithStatus(http.StatusForbidden)
            return
        }
        
        userID, err := strconv.Atoi(userIDStr)
        if err != nil {
            c.AbortWithStatus(http.StatusBadRequest)
            return
        }
        
        targetUser, err = models.GetUserByID(uint(userID))
        if err != nil || targetUser == nil {
            c.AbortWithStatus(http.StatusNotFound)
            return
        }
        
        isOwnProfile = false
    }
    
    // Validation
    errors := make(map[string]string)
    
    // Validate name
    if name == "" {
        errors["name"] = "Name is required"
    } else if len(name) < 2 {
        errors["name"] = "Name must be at least 2 characters long"
    }
    
    // Validate email
    if email == "" {
        errors["email"] = "Email is required"
    } else if !isValidEmail(email) {
        errors["email"] = "Please enter a valid email address"
    } else {
        // Check if email is already taken by another user
        existingUser, _ := models.GetUserByEmail(email)
        if existingUser != nil && existingUser.ID != targetUser.ID {
            errors["email"] = "Email is already taken"
        }
    }
    
    // Validate password change (only if password fields are filled)
    passwordChange := newPassword != "" || confirmPassword != "" || currentPassword != ""
    
    if passwordChange {
        // Only require current password verification for own profile
        if isOwnProfile {
            if currentPassword == "" {
                errors["current_password"] = "Current password is required"
            } else if !internal.CheckPasswordHash(currentPassword, targetUser.Password) {
                errors["current_password"] = "Current password is incorrect"
            }
        }
        
        if newPassword == "" {
            errors["new_password"] = "New password is required"
        } else if len(newPassword) < 6 {
            errors["new_password"] = "Password must be at least 6 characters long"
        }
        
        if confirmPassword == "" {
            errors["confirm_password"] = "Please confirm your new password"
        } else if newPassword != confirmPassword {
            errors["confirm_password"] = "Passwords do not match"
        }
    }
    
    // Validate admin status change (only admins can modify this and not for themselves)
    if !isOwnProfile && isAdmin {
        // Prevent removing admin status from the last admin
        if isAdminStr == "" && targetUser.IsAdmin {
            // Check if this is the last admin
            adminCount, err := models.CountAdminUsers()
            if err != nil {
                errors["general"] = "Failed to validate admin status"
            } else if adminCount <= 1 {
                errors["general"] = "Cannot remove admin privileges - at least one admin must remain"
            }
        }
    }
    
    // If there are validation errors, re-render the form
    if len(errors) > 0 {
        component := auth.ProfilePage(c, targetUser, isOwnProfile, errors)
        component.Render(c.Request.Context(), c.Writer)
        return
    }
    
    // Update user data
    targetUser.Name = name
    targetUser.Email = email
    
    // Update admin status (only if admin is updating another user's profile)
    if !isOwnProfile && isAdmin {
        targetUser.IsAdmin = (isAdminStr == "on")
    }
    
    // Update password if provided
    if passwordChange && newPassword != "" {
        hashedPassword := internal.Encrypt(newPassword)
        targetUser.Password = hashedPassword
    }
    
    // Save to database
    err := models.UpdateUser(targetUser)
    if err != nil {
        errors["general"] = "Failed to update profile. Please try again."
        component := auth.ProfilePage(c, targetUser, isOwnProfile, errors)
        component.Render(c.Request.Context(), c.Writer)
        return
    }
    
    // Update session if user updated their own profile
    if isOwnProfile {
        c.Set("user", targetUser)
    }
    
    // Success - redirect to profile page
    if isOwnProfile {
        c.Redirect(http.StatusSeeOther, "/profile")
    } else {
        c.Redirect(http.StatusSeeOther, "/profile/"+strconv.Itoa(int(targetUser.ID)))
    }
}

// Helper function to validate email format
func isValidEmail(email string) bool {
    // Simple email validation
    return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) > 5
}