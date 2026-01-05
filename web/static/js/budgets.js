// Modal functions
function showBudgetModal() {
    document.getElementById('budgetModal').style.display = 'block';
    document.getElementById('budgetForm').reset();
    document.getElementById('formMessage').style.display = 'none';
}

function hideBudgetModal() {
    document.getElementById('budgetModal').style.display = 'none';
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('budgetModal');
    if (event.target === modal) {
        hideBudgetModal();
    }
}

// Save budget (placeholder for now)
async function saveBudget(event) {
    event.preventDefault();
    
    const formMessage = document.getElementById('formMessage');
    const form = event.target;
    const formData = new FormData(form);
    
    const data = {
        category_id: parseInt(formData.get('category_id')),
        year: parseInt(formData.get('year')),
        amount: parseFloat(formData.get('amount'))
    };
    
    console.log('Budget data to save:', data);
    
    // Design-only demonstration message
    formMessage.className = 'form-message success';
    formMessage.textContent = 'Design Demo: Budget saved successfully (Frontend logic only)!';
    formMessage.style.display = 'block';
    
    setTimeout(() => {
        hideBudgetModal();
    }, 2000);
}
