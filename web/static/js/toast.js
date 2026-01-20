class Toast {
    constructor() {
        this.container = document.createElement('div');
        this.container.className = 'toast-container';
        document.body.appendChild(this.container);
    }

    show(message, type = 'info', title = '') {
        const toast = document.createElement('div');
        toast.className = `toast ${type}`;
        
        const icons = {
            success: '✅',
            error: '❌',
            warning: '⚠️',
            info: 'ℹ️'
        };

        toast.innerHTML = `
            <div class="toast-icon">${icons[type] || icons.info}</div>
            <div class="toast-content">
                ${title ? `<div class="toast-title">${title}</div>` : ''}
                <div class="toast-message">${message}</div>
            </div>
            <button class="toast-close">&times;</button>
        `;

        this.container.appendChild(toast);

        // Trigger animation
        setTimeout(() => toast.classList.add('show'), 10);

        const close = () => {
            toast.classList.remove('show');
            setTimeout(() => toast.remove(), 300);
        };

        toast.querySelector('.toast-close').onclick = close;

        // Auto close after 5 seconds
        setTimeout(close, 5000);
    }

    success(message, title = 'Success') { this.show(message, 'success', title); }
    error(message, title = 'Error') { this.show(message, 'error', title); }
    warning(message, title = 'Warning') { this.show(message, 'warning', title); }
    info(message, title = 'Information') { this.show(message, 'info', title); }
}

const toast = new Toast();

// Automatically handle query parameters for errors
window.addEventListener('load', () => {
    const params = new URLSearchParams(window.location.search);
    if (params.has('error')) {
        const error = params.get('error');
        if (error === 'forbidden') {
            toast.error('You do not have permission to access that page.', 'Access Denied');
        } else if (error === 'unauthorized') {
            toast.warning('Please log in to continue.', 'Authentication Required');
        }
    }
});
