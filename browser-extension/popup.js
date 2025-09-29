// BriefBot Extension Popup JavaScript - Shadcn Aesthetic Edition

class BriefBotExtension {
    constructor() {
        this.API_BASE_URL = 'http://localhost:8080';
        this.currentTab = null;
        this.init();
    }

    async init() {
        await this.loadCurrentTab();
        await this.loadUserId();
        this.setupEventListeners();
        this.checkBackendStatus();
        this.updateUI();
        this.addAestheticEnhancements();
    }

    addAestheticEnhancements() {
        // Add subtle hover effects to buttons
        const buttons = document.querySelectorAll('.button');
        buttons.forEach(button => {
            button.addEventListener('mouseenter', () => {
                button.style.transform = 'translateY(-1px)';
            });
            
            button.addEventListener('mouseleave', () => {
                button.style.transform = 'translateY(0)';
            });
        });

        // Add focus effects to inputs
        const inputs = document.querySelectorAll('.input');
        inputs.forEach(input => {
            input.addEventListener('focus', () => {
                input.parentElement?.classList.add('input-focused');
            });
            
            input.addEventListener('blur', () => {
                input.parentElement?.classList.remove('input-focused');
            });
        });
    }

    async loadCurrentTab() {
        try {
            const tabs = await chrome.tabs.query({ active: true, currentWindow: true });
            this.currentTab = tabs[0];
            
            if (this.currentTab) {
                const pageUrl = document.getElementById('pageUrl');
                const pageTitle = document.getElementById('pageTitle');
                
                pageUrl.textContent = this.truncateUrl(this.currentTab.url, 50);
                pageTitle.textContent = this.currentTab.title || 'Untitled Page';
                
                // Add subtle animation
                pageUrl.style.animation = 'fadeIn 0.3s ease';
                pageTitle.style.animation = 'fadeIn 0.3s ease 0.1s both';
            }
        } catch (error) {
            console.error('Error loading current tab:', error);
            this.showStatus('Error loading current tab information', 'error');
        }
    }

    truncateUrl(url, maxLength) {
        if (url.length <= maxLength) return url;
        return url.substring(0, maxLength - 3) + '...';
    }

    async loadUserId() {
        try {
            const result = await chrome.storage.local.get(['userId']);
            if (result.userId) {
                const userIdInput = document.getElementById('userId');
                userIdInput.value = result.userId;
                
                // Add success animation
                userIdInput.style.animation = 'fadeIn 0.3s ease';
                
                // Show subtle success indicator
                this.showStatus('User ID loaded', 'info');
                setTimeout(() => this.hideStatus(), 1500);
            }
        } catch (error) {
            console.error('Error loading user ID:', error);
        }
    }

    setupEventListeners() {
        // Save user ID with enhanced UX
        document.getElementById('saveUserId').addEventListener('click', async () => {
            const userId = document.getElementById('userId').value;
            const saveBtn = document.getElementById('saveUserId');
            
            if (userId && userId > 0) {
                try {
                    // Show loading state
                    saveBtn.textContent = 'Saving...';
                    saveBtn.disabled = true;
                    
                    await chrome.storage.local.set({ userId: parseInt(userId) });
                    
                    // Show success with animation
                    saveBtn.textContent = 'Saved!';
                    saveBtn.classList.add('button-secondary');
                    saveBtn.classList.remove('button-secondary');
                    
                    this.showStatus('User ID saved successfully!', 'success');
                    
                    // Reset button after success
                    setTimeout(() => {
                        saveBtn.textContent = 'Save';
                        saveBtn.disabled = false;
                        this.hideStatus();
                    }, 2000);
                    
                } catch (error) {
                    console.error('Error saving user ID:', error);
                    this.showStatus('Error saving user ID', 'error');
                    saveBtn.textContent = 'Save';
                    saveBtn.disabled = false;
                }
            } else {
                this.showStatus('Please enter a valid user ID', 'error');
            }
        });

        // Submit URL with enhanced UX
        document.getElementById('submitBtn').addEventListener('click', () => {
            this.submitUrl();
        });

        // Enter key to submit
        document.getElementById('userId').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                document.getElementById('saveUserId').click();
            }
        });
    }

    async checkBackendStatus() {
        try {
            const response = await fetch(`${this.API_BASE_URL}/users`, {
                method: 'GET',
                timeout: 5000
            });
            
            const backendStatus = document.getElementById('backendStatus');
            const backendIndicator = document.getElementById('backendIndicator');
            
            if (response.ok) {
                backendStatus.textContent = 'Online';
                backendIndicator.className = 'status-indicator online';
            } else {
                backendStatus.textContent = 'Offline';
                backendIndicator.className = 'status-indicator offline';
            }
        } catch (error) {
            const backendStatus = document.getElementById('backendStatus');
            const backendIndicator = document.getElementById('backendIndicator');
            backendStatus.textContent = 'Offline';
            backendIndicator.className = 'status-indicator offline';
            console.error('Backend check failed:', error);
        }
    }

    updateUI() {
        const userId = document.getElementById('userId').value;
        const submitBtn = document.getElementById('submitBtn');
        
        if (!userId || userId <= 0) {
            submitBtn.disabled = true;
            this.showStatus('Please enter and save your User ID first', 'info');
        } else if (!this.currentTab || !this.currentTab.url) {
            submitBtn.disabled = true;
            this.showStatus('Unable to get current page information', 'error');
        } else {
            submitBtn.disabled = false;
            this.hideStatus();
        }
    }

    async submitUrl() {
        const userId = parseInt(document.getElementById('userId').value);
        const url = this.currentTab.url;
        const title = this.currentTab.title || 'Untitled Page';

        if (!userId || userId <= 0) {
            this.showStatus('Please enter and save your User ID', 'error');
            return;
        }

        if (!url) {
            this.showStatus('No URL to submit', 'error');
            return;
        }

        // Show loading state with enhanced UX
        const submitBtn = document.getElementById('submitBtn');
        const btnText = submitBtn.querySelector('.button-text');
        const btnIcon = submitBtn.querySelector('.button-icon');
        const spinner = submitBtn.querySelector('.loading-spinner');
        
        submitBtn.disabled = true;
        btnText.textContent = 'Adding...';
        btnIcon.style.display = 'none';
        spinner.style.display = 'flex';
        
        // Add loading animation to button
        submitBtn.classList.add('button-loading');

        try {
            const response = await fetch(`${this.API_BASE_URL}/items`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    user_id: userId,
                    url: url
                })
            });

            const result = await response.json();
            
            if (response.ok) {
                // Success with enhanced animation
                this.showStatus(`✅ "${this.truncateText(title, 30)}" added to your reading list!`, 'success');
                
                // Add success animation to the info card
                const infoCard = document.querySelector('.info-card');
                infoCard.style.animation = 'successPulse 0.6s ease';
                
                // Reset button after success with animation
                setTimeout(() => {
                    submitBtn.disabled = false;
                    btnText.textContent = 'Add to BriefBot';
                    btnIcon.style.display = 'flex';
                    spinner.style.display = 'none';
                    submitBtn.classList.remove('button-loading');
                    infoCard.style.animation = '';
                    this.hideStatus();
                }, 3000);

            } else {
                throw new Error(result.error || 'Failed to submit URL');
            }

        } catch (error) {
            console.error('Submission error:', error);
            this.showStatus(`❌ Error: ${error.message}`, 'error');
            
            // Reset button after error
            submitBtn.disabled = false;
            btnText.textContent = 'Add to BriefBot';
            btnIcon.style.display = 'flex';
            spinner.style.display = 'none';
            submitBtn.classList.remove('button-loading');
        }
    }

    truncateText(text, maxLength) {
        if (text.length <= maxLength) return text;
        return text.substring(0, maxLength - 3) + '...';
    }

    showStatus(message, type) {
        const statusDiv = document.getElementById('statusMessage');
        const statusText = document.getElementById('statusText');
        
        statusText.textContent = message;
        statusDiv.className = `status-message ${type}`;
        statusDiv.style.display = 'flex';
        
        // Add entrance animation
        statusDiv.style.animation = 'slideIn 0.2s ease';
    }

    hideStatus() {
        const statusDiv = document.getElementById('statusMessage');
        statusDiv.style.display = 'none';
    }
}

// Add CSS animations
const style = document.createElement('style');
style.textContent = `
    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(4px); }
        to { opacity: 1; transform: translateY(0); }
    }
    
    @keyframes successPulse {
        0% { transform: scale(1); }
        50% { transform: scale(1.02); }
        100% { transform: scale(1); }
    }
    
    .button-loading {
        opacity: 0.8;
    }
    
    .input-focused {
        transform: translateY(-1px);
    }
`;
document.head.appendChild(style);

// Initialize the extension when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new BriefBotExtension();
});