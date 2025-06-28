import { initializeThemeSwitcher } from './theme.js';
import { elements, translations, resetUi, displayInitialState } from './ui.js';
import { submitUrl, connectToJobEvents } from './api.js';

document.addEventListener('DOMContentLoaded', () => {
    initializeThemeSwitcher();

    elements.pasteButton.addEventListener('click', async () => {
        try {
            elements.urlInput.value = await navigator.clipboard.readText();
        } catch (err) {
            console.error('Failed to read clipboard contents: ', err);
        }
    });

    elements.form.addEventListener('submit', async (event) => {
        event.preventDefault();

        elements.form.style.display = 'none';
        elements.statusArea.style.display = '';
        elements.resetButton.style.display = '';
        elements.spinner.style.display = '';
        elements.statusText.textContent = translations.submittingText;

        try {
            const initialData = await submitUrl(elements.form);
            const initialState = { ...initialData, status: 'processing' };

            sessionStorage.setItem('lastJob', JSON.stringify(initialState));
            displayInitialState(initialState, connectToJobEvents);
        } catch (error) {
            console.error('Submit Error:', error);
            elements.statusText.textContent = 'Error';
            elements.errorMessage.textContent = error.message;
            elements.spinner.style.display = 'none';
        }
    });

    elements.resetButton.addEventListener('click', () => {
        sessionStorage.removeItem('lastJob');
        resetUi();
    });

    const lastJobJSON = sessionStorage.getItem('lastJob');
    if (lastJobJSON) {
        const lastJob = JSON.parse(lastJobJSON);
        displayInitialState(lastJob, connectToJobEvents);
    } else {
        resetUi();
    }
});

