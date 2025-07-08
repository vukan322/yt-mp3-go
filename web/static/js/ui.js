export const elements = {
    body: document.body,
    infoForm: document.getElementById('info-form'),
    urlInput: document.getElementById('url-input'),
    pasteButton: document.getElementById('paste-button'),
    submitButton: document.getElementById('submit-button'),
    statusArea: document.getElementById('status-area'),
    resetButton: document.getElementById('reset-button'),
    thumbnail: document.getElementById('thumbnail'),
    videoTitle: document.getElementById('video-title'),
    qualitySelectorArea: document.getElementById('quality-selector-area'),
    qualitySelector: document.querySelector('.quality-selector'),
    qualitySlider: document.querySelector('.quality-slider'),
    qualityOptions: document.querySelectorAll('.quality-option'),
    qualityDescription: document.getElementById('quality-description'),
    startDownloadButton: document.getElementById('start-download-button'),
    spinner: document.getElementById('spinner'),
    statusText: document.getElementById('status-text'),
    downloadLink: document.getElementById('download-link'),
    errorMessage: document.getElementById('error-message'),
};

export const translations = elements.statusArea.dataset;
export const basePath = elements.body.dataset.basePath || '';

export function updateSliderPosition(targetButton) {
    if (!targetButton || !elements.qualitySlider) return;
    const { offsetLeft, offsetWidth } = targetButton;
    elements.qualitySlider.style.left = `${offsetLeft}px`;
    elements.qualitySlider.style.width = `${offsetWidth}px`;
}

export function resetUi() {
    elements.infoForm.style.display = 'flex';
    elements.statusArea.style.display = 'none';
    elements.statusArea.classList.remove('is-processing');
    elements.urlInput.value = ''; // Bug 2 Fixed
}

export function showSubmittingState() {
    elements.infoForm.style.display = 'none';
    elements.statusArea.style.display = 'block';
    elements.statusText.textContent = translations.submittingText;
    elements.spinner.style.display = 'block';
    elements.statusArea.classList.add('is-processing');
    elements.qualitySelectorArea.style.display = 'none';
    elements.startDownloadButton.style.display = 'none';
    elements.resetButton.style.display = 'none';
}

export function showInfoResult(metadata) {
    elements.statusArea.classList.remove('is-processing');
    elements.spinner.style.display = 'none';
    elements.statusText.textContent = '';
    elements.resetButton.style.display = 'block';

    elements.thumbnail.src = metadata.thumbnail;
    elements.thumbnail.style.display = 'block';
    elements.videoTitle.textContent = metadata.title;
    elements.qualitySelectorArea.style.display = 'block';
    elements.startDownloadButton.style.display = 'block';

    setTimeout(() => updateSliderPosition(document.querySelector('.quality-option.active')), 10);
}

export function showDownloadInProgress() {
    elements.qualitySelectorArea.style.display = 'none';
    elements.startDownloadButton.style.display = 'none';
    elements.spinner.style.display = 'block';
    elements.statusText.textContent = translations.processingText;
    elements.statusArea.classList.add('is-processing');
}

export function showDownloadResult(job, state) {
    elements.statusArea.classList.remove('is-processing');
    elements.spinner.style.display = 'none';
    elements.resetButton.style.display = 'block';

    if (job.status === 'complete') {
        elements.statusText.textContent = translations.readyText;
        elements.downloadLink.href = `${basePath}/${job.filePath}`;
        elements.downloadLink.click();

        setTimeout(() => {
            showInfoResult({ thumbnail: state.thumbnail, title: state.title });
        }, 2500);

    } else if (job.status === 'failed') {
        elements.statusText.textContent = translations.failedText;
        elements.errorMessage.textContent = job.error;
        elements.errorMessage.style.display = 'block';
    }
}

export function showError(message) {
    elements.statusArea.classList.remove('is-processing');
    elements.spinner.style.display = 'none';
    elements.qualitySelectorArea.style.display = 'none';
    elements.startDownloadButton.style.display = 'none';

    elements.statusText.textContent = translations.failedText;
    elements.errorMessage.textContent = message;
    elements.errorMessage.style.display = 'block';
    elements.resetButton.style.display = 'block';
}

export function restoreUiFromState(state) {
    if (!state || !state.videoID) return false;

    elements.infoForm.style.display = 'none';
    elements.statusArea.style.display = 'block';
    elements.thumbnail.src = state.thumbnail;
    elements.thumbnail.style.display = 'block';
    elements.videoTitle.textContent = state.title;

    if (state.jobID && state.jobStatus !== 'complete' && state.jobStatus !== 'failed') {
        showDownloadInProgress();
    } else if (state.jobStatus) {
        const job = { status: state.jobStatus, filePath: state.filePath, error: state.error };
        showDownloadResult(job, state);
    } else {
        showInfoResult({ thumbnail: state.thumbnail, title: state.title });
    }
    return true;
}
