// BriefBot Extension Background Service Worker - Enhanced Edition

// Extension installation and update handling
chrome.runtime.onInstalled.addListener((details) => {
    console.log('BriefBot Extension installed/updated:', details.reason);
    
    // Set default settings with aesthetic preferences
    chrome.storage.local.set({
        apiUrl: 'http://localhost:8080',
        autoSubmit: false,
        showNotifications: true,
        theme: 'system' // system, light, dark
    });

    // Show welcome notification with elegant styling
    if (details.reason === 'install') {
        chrome.notifications.create('welcome', {
            type: 'basic',
            iconUrl: 'icons/icon128.png',
            title: 'BriefBot Extension',
            message: 'Your elegant reading companion is ready! Click the extension icon to start saving articles.',
            contextMessage: 'Set your User ID in the popup to get started.'
        });
    }
});

// Handle messages from content scripts and popup
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    console.log('Background received message:', request);

    switch (request.action) {
        case 'getTabInfo':
            getCurrentTabInfo().then(sendResponse);
            return true; // Keep message channel open for async response

        case 'checkBackend':
            checkBackendStatus().then(sendResponse);
            return true;

        case 'submitUrl':
            submitUrlToBackend(request.data).then(sendResponse);
            return true;

        default:
            console.log('Unknown action:', request.action);
    }
});

// Get current tab information with enhanced details
async function getCurrentTabInfo() {
    try {
        const tabs = await chrome.tabs.query({ active: true, currentWindow: true });
        if (tabs[0]) {
            return {
                success: true,
                tab: {
                    url: tabs[0].url,
                    title: tabs[0].title,
                    id: tabs[0].id,
                    favIconUrl: tabs[0].favIconUrl
                }
            };
        }
        return { success: false, error: 'No active tab found' };
    } catch (error) {
        console.error('Error getting tab info:', error);
        return { success: false, error: error.message };
    }
}

// Check if backend is accessible with timeout
async function checkBackendStatus() {
    try {
        const storage = await chrome.storage.local.get(['apiUrl']);
        const apiUrl = storage.apiUrl || 'http://localhost:8080';
        
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 5000);
        
        const response = await fetch(`${apiUrl}/users`, {
            method: 'GET',
            signal: controller.signal
        });
        
        clearTimeout(timeoutId);
        
        return {
            success: response.ok,
            status: response.status,
            online: response.ok
        };
    } catch (error) {
        console.error('Backend check failed:', error);
        return {
            success: false,
            online: false,
            error: error.message
        };
    }
}

// Submit URL to backend with enhanced error handling
async function submitUrlToBackend(data) {
    try {
        const storage = await chrome.storage.local.get(['apiUrl', 'showNotifications']);
        const apiUrl = storage.apiUrl || 'http://localhost:8080';
        const showNotifications = storage.showNotifications !== false;
        
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 30000); // 30 second timeout
        
        const response = await fetch(`${apiUrl}/items`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                user_id: data.user_id,
                url: data.url
            }),
            signal: controller.signal
        });
        
        clearTimeout(timeoutId);
        const result = await response.json();
        
        if (response.ok) {
            // Show success notification if enabled
            if (showNotifications) {
                chrome.notifications.create('submit-success', {
                    type: 'basic',
                    iconUrl: 'icons/icon128.png',
                    title: '✅ Added to BriefBot',
                    message: `"${data.title || 'Page'}" has been saved to your reading list.`,
                    contextMessage: 'Processing in background...'
                });
            }
            
            return {
                success: true,
                data: result
            };
        } else {
            throw new Error(result.error || `HTTP ${response.status}: ${response.statusText}`);
        }
    } catch (error) {
        console.error('URL submission failed:', error);
        
        // Show error notification if enabled
        const storage = await chrome.storage.local.get(['showNotifications']);
        if (storage.showNotifications !== false) {
            let errorMessage = 'Connection failed';
            if (error.name === 'AbortError') {
                errorMessage = 'Request timeout - backend may be slow';
            } else if (error.message.includes('Failed to fetch')) {
                errorMessage = 'Cannot connect to BriefBot backend';
            } else {
                errorMessage = error.message;
            }
            
            chrome.notifications.create('submit-error', {
                type: 'basic',
                iconUrl: 'icons/icon128.png',
                title: '❌ Submission Failed',
                message: errorMessage,
                contextMessage: 'Please check your connection and try again.'
            });
        }
        
        return {
            success: false,
            error: error.message
        };
    }
}

// Handle extension icon click (optional additional functionality)
chrome.action.onClicked.addListener((tab) => {
    // This only fires if the popup is disabled
    console.log('Extension icon clicked on tab:', tab.id);
});

// Context menu for right-click URL submission with elegant design
chrome.runtime.onInstalled.addListener(() => {
    chrome.contextMenus.create({
        id: 'submitToBriefBot',
        title: '📚 Save to BriefBot',
        contexts: ['page', 'link']
    });
});

chrome.contextMenus.onClicked.addListener(async (info, tab) => {
    if (info.menuItemId === 'submitToBriefBot') {
        const url = info.linkUrl || info.pageUrl;
        const title = tab.title || 'Untitled';
        
        // Get user ID from storage
        const storage = await chrome.storage.local.get(['userId']);
        const userId = storage.userId;
        
        if (!userId) {
            chrome.notifications.create('no-userid', {
                type: 'basic',
                iconUrl: 'icons/icon128.png',
                title: 'BriefBot',
                message: 'Please set your User ID in the extension popup first.',
                contextMessage: 'Click the extension icon to configure.'
            });
            return;
        }
        
        // Submit the URL with loading state
        const result = await submitUrlToBackend({
            user_id: userId,
            url: url,
            title: title
        });
        
        if (!result.success) {
            console.error('Context menu submission failed:', result.error);
        }
    }
});

// Optional: Handle tab updates to track page changes
chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
    if (changeInfo.status === 'complete' && tab.active) {
        // Could update badge or other UI elements here
        console.log('Page loaded:', tab.url);
    }
});

// Optional: Handle extension startup
chrome.runtime.onStartup.addListener(() => {
    console.log('BriefBot extension started');
    // Could perform cleanup or initialization tasks here
});