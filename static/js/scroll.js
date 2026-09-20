const observer = new IntersectionObserver(
  (entries) => {
    entries.forEach((entry) => {
      const card = entry.target.querySelector(".target-element");
      if (!card) return;
      if (entry.isIntersecting) {
        card.classList.add("show");
      } else {
        card.classList.remove("show");
      }
    });
  },
  { threshold: 0.1 }
);

const targetElements = document.querySelectorAll(".target-element");
targetElements.forEach((el) => observer.observe(el.parentElement));