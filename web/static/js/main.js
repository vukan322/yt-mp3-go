import { initializeThemeSwitcher } from './theme.js';
import { elements, resetUi, showSubmittingState, showInfoResult, showDownloadInProgress, showDownloadResult, showError, restoreUiFromState, updateSliderPosition } from './ui.js';
import { getInfo, startDownload, connectToJobEvents } from './api.js';

document.addEventListener('DOMContentLoaded', () => {
    initializeThemeSwitcher();
    let state = JSON.parse(sessionStorage.getItem('yt-downloader-state')) || {};

    elements.pasteButton.addEventListener('click', async () => {
        try {
            elements.urlInput.value = await navigator.clipboard.readText();
        } catch (err) {
            console.error('Failed to read clipboard contents: ', err);
        }
    });

    elements.infoForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        showSubmittingState();
        elements.submitButton.disabled = true;

        try {
            const metadata = await getInfo(elements.infoForm);
            state = { videoID: metadata.id, title: metadata.title, thumbnail: metadata.thumbnail };
            sessionStorage.setItem('yt-downloader-state', JSON.stringify(state));
            showInfoResult(metadata);
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
        showDownloadInProgress();

        try {
            const { jobID } = await startDownload(state.videoID, state.quality || 'high');
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

    elements.resetButton.addEventListener('click', () => {
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
