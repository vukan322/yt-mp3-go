import { basePath } from './ui.js';

export async function getInfo(formElement) {
    const formData = new FormData(formElement);
    const response = await fetch(`${basePath}/info`, {
        method: 'POST',
        body: formData,
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Server responded with an error.');
    }
    return await response.json();
}

export async function startDownload(videoID, quality) {
    const response = await fetch(`${basePath}/download`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ videoID, quality }),
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Server responded with an error.');
    }
    return await response.json();
}

export function connectToJobEvents(jobID, onUpdate) {
    let jobIsFinished = false;
    const eventSource = new EventSource(`${basePath}/events?id=${jobID}`);

    eventSource.onmessage = (e) => {
        const job = JSON.parse(e.data);
        onUpdate(job);

        if (job.status === 'complete' || job.status === 'failed') {
            jobIsFinished = true;
            eventSource.close();
        }
    };

    eventSource.onerror = () => {
        eventSource.close();
        if (!jobIsFinished) {
            onUpdate({ status: 'failed', error: 'Connection to server lost.' });
        }
    };
}
