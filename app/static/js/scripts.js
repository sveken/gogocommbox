// DOM elements
document.addEventListener('DOMContentLoaded', function() {
    // Get all copy buttons
    const copyButtons = document.querySelectorAll('.copy-btn');
    const toast = document.getElementById('toast');
    
    // Add click event to all copy buttons
    copyButtons.forEach(button => {
        button.addEventListener('click', function() {
            // Get URL from data attribute
            const url = this.getAttribute('data-url');
            
            // Copy to clipboard
            navigator.clipboard.writeText(url).then(() => {
                // Show toast notification
                showToast();
            }).catch(err => {
                console.error('Could not copy text: ', err);
            });
        });
    });
    
    // Function to show toast notification
    function showToast() {
        toast.classList.add('show');
        
        // Hide toast after 3 seconds
        setTimeout(() => {
            toast.classList.remove('show');
        }, 3000);
    }
    
    // Add file icons based on extension
    document.querySelectorAll('.file-name i.fa-file').forEach(icon => {
        const filename = icon.nextElementSibling.textContent.trim();
        const extension = filename.split('.').pop().toLowerCase();
        
        // Map of extensions to Font Awesome icons
        const iconMap = {
            'pdf': 'fa-file-pdf',
            'doc': 'fa-file-word',
            'docx': 'fa-file-word',
            'xls': 'fa-file-excel',
            'xlsx': 'fa-file-excel',
            'ppt': 'fa-file-powerpoint',
            'pptx': 'fa-file-powerpoint',
            'jpg': 'fa-file-image',
            'jpeg': 'fa-file-image',
            'png': 'fa-file-image',
            'gif': 'fa-file-image',
            'svg': 'fa-file-image',
            'mp3': 'fa-file-audio',
            'wav': 'fa-file-audio',
            'mp4': 'fa-file-video',
            'mov': 'fa-file-video',
            'zip': 'fa-file-archive',
            'rar': 'fa-file-archive',
            'txt': 'fa-file-alt',
            'js': 'fa-file-code',
            'html': 'fa-file-code',
            'css': 'fa-file-code',
            'json': 'fa-file-code',
            'py': 'fa-file-code',
            'java': 'fa-file-code',
            'c': 'fa-file-code',
            'cpp': 'fa-file-code',
            'cs': 'fa-file-code',
            'go': 'fa-file-code',
        };
        
        if (iconMap[extension]) {
            icon.classList.remove('fa-file');
            icon.classList.add(iconMap[extension]);
        }
    });
});

// Function to format file size
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}