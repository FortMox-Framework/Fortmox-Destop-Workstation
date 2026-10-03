const menuButton = document.querySelector('.docs-menu-button');
const sidebar = document.querySelector('.docs-sidebar');

if (menuButton && sidebar) {
  menuButton.addEventListener('click', () => {
    const open = menuButton.getAttribute('aria-expanded') === 'true';
    menuButton.setAttribute('aria-expanded', String(!open));
    sidebar.classList.toggle('is-open', !open);
  });
}

const search = document.querySelector('#docs-search');
if (search) {
  search.addEventListener('input', () => {
    const query = search.value.trim().toLocaleLowerCase();
    document.querySelectorAll('.doc-nav-group').forEach((group) => {
      let visible = false;
      group.querySelectorAll('li').forEach((item) => {
        const match = item.textContent.toLocaleLowerCase().includes(query);
        item.hidden = !match;
        visible ||= match;
      });
      group.hidden = !visible;
    });
  });
}
