document.addEventListener('DOMContentLoaded', () => {
  const sections = document.querySelectorAll('#index-page section, #index-page .hero-img');

  sections.forEach((section) => {
    section.style.transition = 'transform 0.3s ease, box-shadow 0.3s ease';
    section.style.willChange = 'transform';

    section.addEventListener('mouseenter', () => {
      section.style.transform = 'translateY(-8px) scale(1.20)';
      section.style.boxShadow = '0 12px 28px rgba(0, 0, 0, 0.45)';
    });

    section.addEventListener('mouseleave', () => {
      section.style.transform = 'translateY(0) scale(1)';
      section.style.boxShadow = 'none';
    });
  });
});