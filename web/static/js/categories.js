// Modal functions
function showCreateModal() {
    document.getElementById('modalTitle').textContent = 'Create New Category';
    document.getElementById('submitBtn').textContent = 'Create Category';
    document.getElementById('categoryId').value = '';
    document.getElementById('createModal').style.display = 'block';
    document.getElementById('createCategoryForm').reset();
    document.getElementById('formMessage').style.display = 'none';
}

function showEditModal(id, name, isActive) {
    document.getElementById('modalTitle').textContent = 'Edit Category';
    document.getElementById('submitBtn').textContent = 'Save Changes';
    document.getElementById('categoryId').value = id;
    document.getElementById('categoryName').value = name;
    document.getElementById('isActive').checked = isActive;
    document.getElementById('createModal').style.display = 'block';
    document.getElementById('formMessage').style.display = 'none';
}

function hideCreateModal() {
    document.getElementById('createModal').style.display = 'none';
}

// Handle Filter Toggle (Instant UI Update)
function handleFilterToggle(checkbox) {
    const activeOnly = checkbox.checked;
    const tableBody = document.getElementById('categoriesTableBody');
    const rows = tableBody.querySelectorAll('tr[data-active]');
    
    rows.forEach(row => {
        const isActive = row.getAttribute('data-active') === 'true';
        if (activeOnly && !isActive) {
            row.style.display = 'none';
        } else {
            row.style.display = '';
        }
    });
}

// Handle Submit (Create or Update)
async function handleCategorySubmit(event) {
    event.preventDefault();
    
    const formMessage = document.getElementById('formMessage');
    const form = event.target;
    const formData = new FormData(form);
    const id = formData.get('id');
    
    const data = {
        name: formData.get('name'),
        is_active: formData.get('is_active') === 'on'
    };
    
    const url = id ? `/api/categories/${id}` : '/api/categories';
    const method = id ? 'PUT' : 'POST';
    
    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (response.ok) {
            formMessage.className = 'form-message success';
            formMessage.textContent = result.message || `Category ${id ? 'updated' : 'created'} successfully!`;
            formMessage.style.display = 'block';
            
            setTimeout(() => {
                hideCreateModal();
                window.location.reload();
            }, 1000);
        } else {
            formMessage.className = 'form-message error';
            formMessage.textContent = result.message || 'Failed to process category';
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
    if (!tableBody) return;
    
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
    
    const activeOnly = document.getElementById('activeOnlyToggle').checked;
    
    tableBody.innerHTML = categories.map(category => {
        const date = new Date(category.created_at).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: '2-digit'
        });

        const statusIcon = '<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 256 256"><path fill="currentColor" d="M144 128a16 16 0 1 1-16-16a16 16 0 0 1 16 16m-64-16a16 16 0 1 0 16 16a16 16 0 0 0-16-16m128 0a16 16 0 1 0 16 16a16 16 0 0 0-16-16"/></svg>';

        return `
            <tr class="${!category.is_active ? 'row-inactive' : ''}" data-id="${category.id}" data-active="${category.is_active}" style="${activeOnly && !category.is_active ? 'display: none;' : ''}">
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
                        <button class="btn-icon" onclick="showEditModal(${category.id}, '${category.name.replace(/'/g, "\\'")}', ${category.is_active})" aria-label="Edit Category" title="Edit Category">
                            <svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 256 256"><path fill="currentColor" d="M227.31 73.37L182.63 28.7a16 16 0 0 0-22.63 0L36.69 152A15.86 15.86 0 0 0 32 163.31V208a16 16 0 0 0 16 16h44.69a15.86 15.86 0 0 0 11.31-4.69L227.31 96a16 16 0 0 0 0-22.63M92.69 208H48v-44.69l88-88L180.69 120ZM192 108.69L147.31 64l24-24L216 84.69Z"/></svg>
                        </button>
                    </div>
                </td>
            </tr>
        `;
    }).join('');
}

// Toggle status
async function toggleStatus(id, currentStatus) {
    if (!confirm(`Are you sure you want to ${currentStatus ? 'deactivate' : 'activate'} this category?`)) {
        return;
    }

    try {
        const response = await fetch(`/api/categories/${id}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
            }
        });
        
        const result = await response.json();
        
        if (response.ok && result.success) {
            // Re-fetch categories to update the list
            fetchCategories();
        } else {
            alert(result.message || 'Failed to toggle status');
        }
    } catch (error) {
        console.error('Error toggling status:', error);
        alert('An error occurred. Please try again.');
    }
}


// Utility function to escape HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
