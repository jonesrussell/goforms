document.getElementById('signupForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Prevent the default form submission

    // Gather form data
    const formData = new FormData(this);
    const data = Object.fromEntries(formData.entries());

    // Send data to the API
    fetch('/api/v1/auth/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    })
    .then(response => response.json())
    .then(data => {
        // Handle success or error response
        console.log('Success:', data);
        // You can also update the UI or show a message to the user
    })
    .catch((error) => {
        console.error('Error:', error);
    });
});