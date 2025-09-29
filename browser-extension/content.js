// BriefBot Extension Content Script
// This script runs on all pages and can extract additional information if needed

(function() {
    'use strict';

    // Extract page metadata
    function extractPageInfo() {
        const info = {
            url: window.location.href,
            title: document.title,
            description: '',
            author: '',
            publishedDate: '',
            keywords: []
        };

        // Try to get meta description
        const descriptionMeta = document.querySelector('meta[name="description"]') || 
                               document.querySelector('meta[property="og:description"]') ||
                               document.querySelector('meta[name="twitter:description"]');
        if (descriptionMeta) {
            info.description = descriptionMeta.getAttribute('content');
        }

        // Try to get author
        const authorMeta = document.querySelector('meta[name="author"]') ||
                          document.querySelector('meta[property="og:article:author"]') ||
                          document.querySelector('meta[name="twitter:creator"]');
        if (authorMeta) {
            info.author = authorMeta.getAttribute('content');
        }

        // Try to get published date
        const dateMeta = document.querySelector('meta[property="article:published_time"]') ||
                        document.querySelector('meta[name="date"]') ||
                        document.querySelector('meta[property="og:article:published_time"]');
        if (dateMeta) {
            info.publishedDate = dateMeta.getAttribute('content');
        }

        // Try to get keywords
        const keywordsMeta = document.querySelector('meta[name="keywords"]');
        if (keywordsMeta) {
            info.keywords = keywordsMeta.getAttribute('content')
                .split(',')
                .map(keyword => keyword.trim())
                .filter(keyword => keyword.length > 0);
        }

        return info;
    }

    // Listen for messages from popup
    chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
        if (request.action === 'getPageInfo') {
            const pageInfo = extractPageInfo();
            sendResponse(pageInfo);
        }
        return true; // Keep the message channel open for async response
    });

    // Optional: Add a small indicator that the extension is active
    // This can be useful for debugging
    console.log('BriefBot extension content script loaded on:', window.location.href);

})();