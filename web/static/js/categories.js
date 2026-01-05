// Modal functions
function showCreateModal() {
    document.getElementById('createModal').style.display = 'block';
    document.getElementById('createCategoryForm').reset();
    document.getElementById('formMessage').style.display = 'none';
}

function hideCreateModal() {
    document.getElementById('createModal').style.display = 'none';
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('createModal');
    if (event.target === modal) {
        hideCreateModal();
    }
}

// Create category
async function createCategory(event) {
    event.preventDefault();
    
    const formMessage = document.getElementById('formMessage');
    const form = event.target;
    const formData = new FormData(form);
    
    const data = {
        name: formData.get('name'),
        is_active: formData.get('is_active') === 'on'
    };
    
    try {
        const response = await fetch('/api/categories', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (response.ok) {
            formMessage.className = 'form-message success';
            formMessage.textContent = result.message || 'Category created successfully!';
            formMessage.style.display = 'block';
            
            // Reload the page after a short delay
            setTimeout(() => {
                window.location.reload();
            }, 1500);
        } else {
            formMessage.className = 'form-message error';
            formMessage.textContent = result.message || 'Failed to create category';
            formMessage.style.display = 'block';
        }
    } catch (error) {
        formMessage.className = 'form-message error';
        formMessage.textContent = 'An error occurred. Please try again.';
        formMessage.style.display = 'block';
        console.error('Error:', error);
    }
}

// Fetch and display categories (for dynamic updates)
async function fetchCategories() {
    try {
        const response = await fetch('/api/categories');
        const result = await response.json();
        
        if (response.ok && result.success) {
            displayCategories(result.data);
        }
    } catch (error) {
        console.error('Error fetching categories:', error);
    }
}

function displayCategories(categories) {
    const tableBody = document.getElementById('categoriesTableBody');
    
    if (!categories || categories.length === 0) {
        tableBody.innerHTML = `
            <tr>
                <td colspan="5" class="empty-state">
                    No categories found. Create your first category!
                </td>
            </tr>
        `;
        return;
    }
    
    tableBody.innerHTML = categories.map(category => {
        const date = new Date(category.created_at).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: '2-digit'
        });

        const statusIcon = category.is_active 
            ? '<svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 256 256"><path fill="currentColor" d="M216 48H40a16 16 0 0 0-16 16v128a16 16 0 0 0 16 16h176a16 16 0 0 0 16-16V64a16 16 0 0 0-16-16m0 144H40V64h176zM112 128a12 12 0 1 1-12-12a12 12 0 0 1 12 12m48 0a12 12 0 1 1-12-12a12 12 0 0 1 12 12"/></svg>'
            : '<svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 256 256"><path fill="currentColor" d="M128 24a104 104 0 1 0 104 104A104.11 104.11 0 0 0 128 24m0 192a88 88 0 1 1 88-88a88.1 88.1 0 0 1-88 88m45.66-93.66a8 8 0 0 1 0 11.32l-32 32a8 8 0 0 1-11.32 0l-32-32a8 8 0 0 1 11.32-11.32L120 132.69V88a8 8 0 0 1 16 0v44.69l11.34-11.35a8 8 0 0 1 11.32 0"/></svg>';

        return `
            <tr class="${!category.is_active ? 'row-inactive' : ''}" data-id="${category.id}">
                <td>#${category.id}</td>
                <td class="font-bold">${escapeHtml(category.name)}</td>
                <td>
                    <span class="status-badge ${category.is_active ? 'active' : 'inactive'}">
                        ${category.is_active ? 'Active' : 'Inactive'}
                    </span>
                </td>
                <td class="text-secondary">${date}</td>
                <td class="text-right">
                    <div class="action-group">
                        <button class="btn-icon" onclick="toggleStatus(${category.id}, ${category.is_active})" 
                                aria-label="${category.is_active ? 'Mark Inactive' : 'Mark Active'}"
                                title="${category.is_active ? 'Mark Inactive' : 'Mark Active'}">
                            ${statusIcon}
                        </button>
                        <button class="btn-icon" aria-label="Edit Category" title="Edit Category">
                            <svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 256 256"><path fill="currentColor" d="M227.31 73.37L182.63 28.7a16 16 0 0 0-22.63 0L36.69 152A15.86 15.86 0 0 0 32 163.31V208a16 16 0 0 0 16 16h44.69a15.86 15.86 0 0 0 11.31-4.69L227.31 96a16 16 0 0 0 0-22.63M92.69 208H48v-44.69l88-88L180.69 120ZM192 108.69L147.31 64l24-24L216 84.69Z"/></svg>
                        </button>
                    </div>
                </td>
            </tr>
        `;
    }).join('');
}

// Toggle status placeholder
function toggleStatus(id, currentStatus) {
    console.log(`Toggle status for category ${id}. Current status: ${currentStatus}`);
    // Future implementation for status toggle goes here
    alert(`Status toggle clicked for Category #${id}. This feature will be implemented soon!`);
}


// Utility function to escape HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
