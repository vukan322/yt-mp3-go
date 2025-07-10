export const elements = {};

export let translations = {};
export let basePath = '';

export function cacheDOMElements() {
    elements.body = document.body;
    elements.infoForm = document.getElementById('info-form');
    elements.urlInput = document.getElementById('url-input');
    elements.pasteButton = document.getElementById('paste-button');
    elements.submitButton = document.getElementById('submit-button');
    elements.statusArea = document.getElementById('status-area');
    elements.resetButton = document.getElementById('reset-button');
    elements.thumbnail = document.getElementById('thumbnail');
    elements.videoTitle = document.getElementById('video-title');
    elements.qualitySelectorArea = document.getElementById('quality-selector-area');
    elements.qualitySelector = document.querySelector('.quality-selector');
    elements.qualitySlider = document.querySelector('.quality-slider');
    elements.qualityOptions = document.querySelectorAll('.quality-option');
    elements.qualityDescription = document.getElementById('quality-description');
    elements.startDownloadButton = document.getElementById('start-download-button');
    elements.spinner = document.getElementById('spinner');
    elements.statusText = document.getElementById('status-text');
    elements.downloadLink = document.getElementById('download-link');
    elements.errorMessage = document.getElementById('error-message');
    elements.downloadOptionsArea = document.getElementById('download-options-area');
    elements.filenameInput = document.getElementById('filename-input');
    elements.filenameResetButton = document.getElementById('filename-reset-button');

    translations = elements.statusArea.dataset;
    basePath = elements.body.dataset.basePath || '';
}

export function updateSliderPosition(targetButton) {
    if (!targetButton || !elements.qualitySlider) return;
    const { offsetLeft, offsetWidth } = targetButton;
    elements.qualitySlider.style.left = `${offsetLeft}px`;
    elements.qualitySlider.style.width = `${offsetWidth}px`;
}

export function resetUi() {
    elements.infoForm.style.display = 'flex';
    elements.statusArea.style.display = 'none';
    elements.downloadOptionsArea.style.display = 'none';
    elements.urlInput.value = '';
    elements.thumbnail.src = '';
    elements.videoTitle.textContent = '';
    elements.filenameInput.value = '';
    elements.statusText.textContent = '';
    elements.errorMessage.textContent = '';
    elements.statusArea.classList.remove('is-processing');
    elements.filenameInput.classList.remove('error');
    elements.errorMessage.classList.remove('visible');
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
    elements.thumbnail.style.display = 'none';
    elements.videoTitle.textContent = '';
}

export function showInfoResult(metadata, state) {
    elements.statusArea.classList.remove('is-processing');
    elements.spinner.style.display = 'none';
    elements.statusText.textContent = '';
    elements.resetButton.style.display = 'block';
    elements.thumbnail.src = metadata.thumbnail;
    elements.thumbnail.style.display = 'block';
    elements.videoTitle.textContent = metadata.title;
    elements.filenameInput.value = state.filename || metadata.title;
    elements.qualitySelectorArea.style.display = 'block';
    elements.downloadOptionsArea.style.display = 'block';
    elements.startDownloadButton.style.display = 'block';
    setTimeout(() => updateSliderPosition(document.querySelector('.quality-option.active')), 10);
}

export function showDownloadInProgress() {
    elements.qualitySelectorArea.style.display = 'none';
    elements.downloadOptionsArea.style.display = 'none';
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

        state.jobID = null;
        state.jobStatus = null;
        state.filePath = null;
        sessionStorage.setItem('yt-downloader-state', JSON.stringify(state));

        setTimeout(() => {
            showInfoResult({ thumbnail: state.thumbnail, title: state.title }, state);
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
    elements.errorMessage.classList.add('visible');
    elements.resetButton.style.display = 'block';
}

export function restoreUiFromState(state) {
    if (!state || !state.videoID) return false;

    elements.infoForm.style.display = 'none';
    elements.statusArea.style.display = 'block';
    elements.thumbnail.src = state.thumbnail;
    elements.thumbnail.style.display = 'block';
    elements.videoTitle.textContent = state.title;
    elements.filenameInput.value = state.filename || state.title;

    if (state.jobID && state.jobStatus !== 'complete' && state.jobStatus !== 'failed') {
        showDownloadInProgress();
    } else if (state.jobStatus) {
        const job = { status: state.jobStatus, filePath: state.filePath, error: state.error };
        showDownloadResult(job, state);
    } else {
        showInfoResult({ thumbnail: state.thumbnail, title: state.title }, state);
    }

    return true;
}

export function showFilenameError() {
    elements.filenameInput.classList.add('error');
    elements.errorMessage.textContent = translations.filenameRequiredText;
    elements.errorMessage.classList.add('visible');
    elements.startDownloadButton.disabled = true;
}

export function clearFilenameError() {
    elements.filenameInput.classList.remove('error');
    elements.errorMessage.classList.remove('visible');
    elements.startDownloadButton.disabled = false;
}
