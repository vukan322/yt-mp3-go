export function initializeThemeSwitcher() {
    const themeSwitcher = document.getElementById('theme-switcher');
    if (!themeSwitcher) return;

    const themeIconLight = document.getElementById('theme-icon-light');
    const themeIconDark = document.getElementById('theme-icon-dark');

    const setTheme = (isDark) => {
        document.body.classList.toggle('dark-theme', isDark);
        themeIconLight.style.display = isDark ? 'block' : 'none';
        themeIconDark.style.display = isDark ? 'none' : 'block';
        localStorage.setItem('theme', isDark ? 'dark' : 'light');
    };

    const isInitialDark = localStorage.getItem('theme') !== 'light';
    setTheme(isInitialDark);

    themeSwitcher.addEventListener('click', () => setTheme(!document.body.classList.contains('dark-theme')));
}

