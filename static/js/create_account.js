// function validateForm() {
//     let isValid = true;

//     // Clear previous errors
//     document.querySelectorAll('.error').forEach(e => e.textContent = '');

//     // Password match validation
//     const password = document.getElementById('password').value;
//     const confirmPassword = document.getElementById('confirmPassword').value;

//     if (password !== confirmPassword) {
//         document.getElementById('confirmPasswordError').textContent =
//             'Passwords do not match';
//         isValid = false;
//     }

//     // Additional validations can be added here

//     return isValid;
// }

// // Real-time password match check
// document.getElementById('confirmPassword').addEventListener('input', function () {
//     const password = document.getElementById('password').value;
//     const confirmError = document.getElementById('confirmPasswordError');

//     if (this.value !== password) {
//         confirmError.textContent = 'Passwords do not match';
//     } else {
//         confirmError.textContent = '';
//     }
// });
