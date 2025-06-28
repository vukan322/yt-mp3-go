import { elements, translations, updateUiForJobStatus } from './ui.js';

export async function submitUrl(formElement) {
    const formData = new FormData(formElement);

    const response = await fetch(`${elements.body.dataset.basePath}/download`, {
        method: 'POST',
        body: formData,
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Server responded with an error.');
    }

    return await response.json();
}

export function connectToJobEvents(job) {
    let jobIsFinished = false;
    const eventSource = new EventSource(`${elements.body.dataset.basePath}/events?id=${job.jobID}`);

    eventSource.onmessage = (e) => {
        const currentJobState = JSON.parse(sessionStorage.getItem('lastJob') || '{}');
        const newUpdate = JSON.parse(e.data);
        const finalState = { ...currentJobState, ...newUpdate };

        if (finalState.status === 'complete' || finalState.status === 'failed') {
            jobIsFinished = true;
            sessionStorage.setItem('lastJob', JSON.stringify(finalState));
            updateUiForJobStatus(finalState);
            eventSource.close();
        }
    };

    eventSource.onerror = () => {
        eventSource.close();
        if (!jobIsFinished) {
            elements.statusText.textContent = translations.failedText;
            elements.errorMessage.textContent = translations.connectionLostText;
            elements.spinner.style.display = 'none';
            elements.resetButton.style.display = '';
        }
    };
}

