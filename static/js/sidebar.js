const opendiv = document.getElementById("sidebar-open");
const openbtn = document.getElementById("sidebar-open-button");
const closebtn = document.getElementById("sidebar-close-button");
const modal = document.getElementById("sidebar-modal");

openbtn.addEventListener("click", () => {
  modal.classList.add("show");
  opendiv.classList.add("hide");
});

closebtn.addEventListener("click", () => {
  modal.classList.remove("show");
  opendiv.classList.remove("hide");
});
