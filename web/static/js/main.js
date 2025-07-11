import { initializeThemeSwitcher } from './theme.js';
import { elements, cacheDOMElements, resetUi, showSubmittingState, showInfoResult, showDownloadInProgress, showDownloadResult, showError, restoreUiFromState, updateSliderPosition, showFilenameError, clearFilenameError, clearUrlError, showUrlError } from './ui.js';
import { getInfo, startDownload, connectToJobEvents, cancelJob } from './api.js';

document.addEventListener('DOMContentLoaded', () => {
    cacheDOMElements();
    initializeThemeSwitcher();

    let state = JSON.parse(sessionStorage.getItem('yt-downloader-state')) || {};
    const youtubeUrlRegex = /^(https?\:\/\/)?(www\.youtube\.com|youtu\.be)\/.+$/;
    let debounceTimer;

    function validateUrl() {
        const isValid = youtubeUrlRegex.test(elements.urlInput.value);
        if (isValid) {
            clearUrlError();
        } else {
            showUrlError();
        }
        return isValid;
    }

    elements.urlInput.addEventListener('input', () => {
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
            validateUrl();
        }, 250);
    });

    elements.pasteButton.addEventListener('click', async () => {
        try {
            const pastedText = await navigator.clipboard.readText();
            elements.urlInput.value = pastedText;
            validateUrl();
        } catch (err) {
            console.error('Failed to read clipboard contents: ', err);
        }
    });

    elements.infoForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        if (!validateUrl()) {
            return;
        }

        showSubmittingState();
        elements.submitButton.disabled = true;

        try {
            const metadata = await getInfo(elements.infoForm);
            state = {
                videoID: metadata.id,
                title: metadata.title,
                thumbnail: metadata.thumbnail,
                filename: metadata.title
            };
            sessionStorage.setItem('yt-downloader-state', JSON.stringify(state));
            showInfoResult(metadata, state);
        } catch (error) {
            console.error('Info Error:', error);
            showError(error.message);
        } finally {
            elements.submitButton.disabled = false;
        }
    });

    elements.qualitySelector.addEventListener('click', (e) => {
        const target = e.target.closest('.quality-option');
        if (target) {
            elements.qualityOptions.forEach(opt => opt.classList.remove('active'));
            target.classList.add('active');
            updateSliderPosition(target);
            state.quality = target.dataset.quality;
            elements.qualityDescription.textContent = target.dataset.desc;
            localStorage.setItem('preferredQuality', state.quality);
        }
    });

    elements.startDownloadButton.addEventListener('click', async () => {
        if (!state.videoID) return;

        const filename = elements.filenameInput.value.trim();
        if (!filename) {
            showFilenameError();
            return;
        }
        state.filename = filename;

        showDownloadInProgress();

        try {
            const { jobID } = await startDownload(state.videoID, state.quality || 'high', state.filename);
            state.jobID = jobID;
            sessionStorage.setItem('yt-downloader-state', JSON.stringify(state));

            connectToJobEvents(jobID, (job) => {
                if (job.status === 'complete' || job.status === 'failed') {
                    state.jobStatus = job.status;
                    state.filePath = job.filePath;
                    state.error = job.error;
                    sessionStorage.setItem('yt-downloader-state', JSON.stringify(state));
                    showDownloadResult(job, state);
                }
            });
        } catch (error) {
            console.error('Download Error:', error);
            showError(error.message);
        }
    });

    elements.filenameInput.addEventListener('input', () => {
        clearFilenameError();

        state.filename = elements.filenameInput.value;
        sessionStorage.setItem('yt-downloader-state', JSON.stringify(state));
    });

    elements.filenameResetButton.addEventListener('click', () => {
        if (state.title) {
            elements.filenameInput.value = state.title;
            state.filename = state.title;
            sessionStorage.setItem('yt-downloader-state', JSON.stringify(state));
        }
    });

    elements.resetButton.addEventListener('click', () => {
        if (state.jobID) {
            cancelJob(state.jobID).catch(err => console.error("Cancel request failed:", err));
        }

        sessionStorage.removeItem('yt-downloader-state');
        state = {};
        resetUi();
    });

    const preferredQuality = localStorage.getItem('preferredQuality') || 'high';
    elements.qualityOptions.forEach(opt => {
        const isActive = opt.dataset.quality === preferredQuality;
        opt.classList.toggle('active', isActive);
        if (isActive) {
            elements.qualityDescription.textContent = opt.dataset.desc;
            state.quality = preferredQuality;
        }
    });

    if (restoreUiFromState(state)) {
        if (state.jobID && state.jobStatus !== 'complete' && state.jobStatus !== 'failed') {
            connectToJobEvents(state.jobID, (job) => {
                if (job.status === 'complete' || job.status === 'failed') {
                    state.jobStatus = job.status;
                    state.filePath = job.filePath;
                    state.error = job.error;
                    sessionStorage.setItem('yt-downloader-state', JSON.stringify(state));
                    showDownloadResult(job, state);
                }
            });
        }
    } else {
        resetUi();
    }
});
