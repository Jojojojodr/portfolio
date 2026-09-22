document.addEventListener("DOMContentLoaded", () => {
  const sections = document.querySelectorAll("#index-page .target-element");
  const imgs = document.querySelectorAll("#index-page .hero-img");

  sections.forEach((section) => {
    section.style.willChange = 'transform';

    section.addEventListener("mouseenter", () => {
      section.style.transform = "translateY(-8px) scale(1.20)";
      section.style.boxShadow = "0 12px 28px rgba(0, 0, 0, 0.45)";
    });

    section.addEventListener("mouseleave", () => {
      section.style.transform = "";
      section.style.boxShadow = "";
    });
  });

  imgs.forEach((img) => {
    img.style.transition = "transform 0.3s ease, box-shadow 0.3s ease";
    img.style.willChange = 'transform';

    img.addEventListener("mouseenter", () => {
      img.style.transform = "translateY(-8px) scale(1.6)";
      img.style.boxShadow = "0 12px 28px rgba(0, 0, 0, 0.45)";
    });

    img.addEventListener("mouseleave", () => {
      img.style.transform = "";
      img.style.boxShadow = "";
    });

  });
});
