// Main JavaScript file for AI in Action app

document.addEventListener('DOMContentLoaded', function() {
    // Initialize Bootstrap tooltips
    const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });

    // Initialize Bootstrap popovers
    const popoverTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="popover"]'));
    popoverTriggerList.map(function (popoverTriggerEl) {
        return new bootstrap.Popover(popoverTriggerEl);
    });

    // Set up HTMX event listeners
    document.body.addEventListener('htmx:afterSwap', function(event) {
        // Reinitialize Bootstrap components after HTMX swaps content
        const newTooltipTriggerList = [].slice.call(event.detail.elt.querySelectorAll('[data-bs-toggle="tooltip"]'));
        newTooltipTriggerList.map(function (tooltipTriggerEl) {
            return new bootstrap.Tooltip(tooltipTriggerEl);
        });

        const newPopoverTriggerList = [].slice.call(event.detail.elt.querySelectorAll('[data-bs-toggle="popover"]'));
        newPopoverTriggerList.map(function (popoverTriggerEl) {
            return new bootstrap.Popover(popoverTriggerEl);
        });
    });
    
    // Global function to close modals after HTMX requests
    window.closeModal = function(modalId) {
        const modalEl = document.getElementById(modalId || 'add-event-modal');
        if (!modalEl) return;
        
        const modalInstance = bootstrap.Modal.getInstance(modalEl);
        if (modalInstance) {
            modalInstance.hide();
        }
    };
    
    // Listen for successful form submissions via HTMX
    document.body.addEventListener('htmx:afterRequest', function(event) {
        // If the request was successful and came from a form inside a modal
        if (event.detail.successful && event.detail.elt.tagName === 'FORM') {
            const modal = event.detail.elt.closest('.modal');
            if (modal) {
                const modalInstance = bootstrap.Modal.getInstance(modal);
                if (modalInstance) {
                    modalInstance.hide();
                }
            }
        }
    });
}); 