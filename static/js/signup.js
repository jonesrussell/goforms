document.getElementById('signupForm').addEventListener('submit', async (event) => {
    event.preventDefault(); // Prevent the default form submission

    const data = gatherFormData(event.target);

    // Validate required fields
    if (!data.email || !data.password) {
        renderMessage('Email and password are required.', 'error');
        return;
    }

    try {
        const responseData = await sendFormData('/api/v1/auth/signup', data);

        console.log('Response Data:', responseData);
        // Check if the response indicates success
        if (responseData.success) {
            renderMessage('Signup successful!', 'success');
        } else {
            renderMessage(responseData.message || 'Signup failed.', 'error');
        }
    } catch (error) {
        console.error('Error:', error);
        renderMessage(error.message, 'error');
    }
});

const gatherFormData = (form) => {
    const formData = new FormData(form);
    return Object.fromEntries(formData.entries());
};

const sendFormData = async (url, data) => {
    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to sign up');
    }

    return await response.json();
};

const renderMessage = (text, type) => {
    const messageArea = document.getElementById('flash-message');
    const messageText = document.getElementById('flash-message-text');

    // Clear previous messages
    messageText.textContent = text;
    messageArea.style.display = 'block'; // Show the message area

    // Add success or error class
    if (type === 'success') {
        messageArea.className = 'flash-message flash-message-success'; // Add success class
    } else {
        messageArea.className = 'flash-message flash-message-error'; // Add error class
    }
};
