export const elements = {
    body: document.body,
    form: document.getElementById('download-form'),
    resetButton: document.getElementById('reset-button'),
    statusArea: document.getElementById('status-area'),
    pasteButton: document.getElementById('paste-button'),
    urlInput: document.getElementById('url-input'),
    spinner: document.getElementById('spinner'),
    statusText: document.getElementById('status-text'),
    downloadLink: document.getElementById('download-link'),
    errorMessage: document.getElementById('error-message'),
    thumbnail: document.getElementById('thumbnail'),
    videoTitle: document.getElementById('video-title'),
    fileSizeText: document.getElementById('file-size-text'),
};

export const translations = elements.statusArea.dataset;
export const basePath = elements.body.dataset.basePath || '';

export function resetUi() {
    elements.statusArea.style.display = 'none';
    elements.form.style.display = '';
    elements.thumbnail.src = '';
    elements.videoTitle.textContent = '';
    elements.statusText.textContent = '';
    elements.errorMessage.textContent = '';
    elements.downloadLink.textContent = '';
    elements.fileSizeText.textContent = '';
    elements.thumbnail.style.display = 'none';
    elements.spinner.style.display = 'none';
    elements.downloadLink.style.display = 'none';
    elements.errorMessage.style.display = 'none';
    elements.fileSizeText.style.display = 'none';
    elements.urlInput.value = '';
}

export function updateUiForJobStatus(job) {
    elements.spinner.style.display = 'none';
    elements.statusArea.classList.remove('processing');
    elements.resetButton.style.display = '';

    if (job.status === 'complete') {
        elements.statusText.textContent = translations.readyText;
        const fileSizeMB = (job.fileSize / (1024 * 1024)).toFixed(2);
        elements.downloadLink.textContent = translations.downloadFileText;
        elements.downloadLink.href = `${basePath}/${job.filePath}`;
        elements.downloadLink.style.display = 'inline-block';
        elements.fileSizeText.textContent = `${translations.fileSizeLabel}: ${fileSizeMB} MB`;
        elements.fileSizeText.style.display = 'block';
    } else if (job.status === 'failed') {
        elements.statusText.textContent = translations.failedText;
        elements.errorMessage.textContent = job.error;
        elements.errorMessage.style.display = 'block';
    }
}

export function displayInitialState(job, onConnect) {
    elements.form.style.display = 'none';
    elements.statusArea.style.display = '';
    elements.resetButton.style.display = '';
    elements.thumbnail.src = job.thumbnail;
    elements.thumbnail.style.display = '';
    elements.videoTitle.textContent = job.title;

    if (job.status === 'pending' || job.status === 'processing') {
        elements.spinner.style.display = '';
        elements.statusText.textContent = translations.processingText;
        elements.statusArea.classList.add('processing');
        onConnect(job);
    } else {
        updateUiForJobStatus(job);
    }
}

