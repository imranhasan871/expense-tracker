// Modal functions
function showExpenseModal() {
    document.getElementById('expenseModal').style.display = 'block';
    document.getElementById('expenseForm').reset();
    document.getElementById('expenseDate').valueAsDate = new Date(); // Reset to today
    document.getElementById('formMessage').style.display = 'none';
}

function hideExpenseModal() {
    document.getElementById('expenseModal').style.display = 'none';
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('expenseModal');
    if (event.target === modal) {
        hideExpenseModal();
    }
}

// Save expense (placeholder for design demo)
async function saveExpense(event) {
    event.preventDefault();
    
    const formMessage = document.getElementById('formMessage');
    const form = event.target;
    const formData = new FormData(form);
    
    const data = {
        amount: parseFloat(formData.get('amount')),
        date: formData.get('date'),
        category_id: parseInt(formData.get('category_id')),
        remarks: formData.get('remarks')
    };
    
    console.log('Expense data to save:', data);
    
    // Design-only demonstration message
    formMessage.className = 'form-message success';
    formMessage.textContent = 'Design Demo: Expense recorded successfully (Frontend only)!';
    formMessage.style.display = 'block';
    
    setTimeout(() => {
        hideExpenseModal();
    }, 2000);
}

// Apply filters (placeholder for design demo)
function applyFilters() {
    const startDate = document.getElementById('filterStartDate').value;
    const endDate = document.getElementById('filterEndDate').value;
    const categoryId = document.getElementById('filterCategory').value;
    
    console.log('Applying filters:', { startDate, endDate, categoryId });
    alert(`Filtering for dates: ${startDate || 'All'} to ${endDate || 'All'}, Category: ${categoryId || 'All'}`);
}
