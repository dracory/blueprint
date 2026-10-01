const { createApp, ref } = Vue;

createApp({
	setup() {
		const firstName = ref('');
		const lastName = ref('');
		const email = ref('');
		const password = ref('');
		const passwordConfirm = ref('');
		const isLoading = ref(false);

		const handleSubmit = async () => {
			if (!firstName.value || !lastName.value || !email.value || !password.value) {
				Notiflix.Notify.failure('Please fill in all fields');
				return;
			}

			if (password.value !== passwordConfirm.value) {
				Notiflix.Notify.failure('Passwords do not match');
				return;
			}

			isLoading.value = true;

			const formData = new FormData();
			formData.append('first_name', firstName.value);
			formData.append('last_name', lastName.value);
			formData.append('email', email.value);
			formData.append('password', password.value);
			formData.append('password_confirm', passwordConfirm.value);
			if (RETURN_URL) {
				formData.append('return', decodeURIComponent(RETURN_URL));
			}

			try {
				const response = await fetch(SUBMIT_AJAX_URL, {
					method: 'POST',
					headers: {
						'Content-Type': 'application/x-www-form-urlencoded',
					},
					body: new URLSearchParams(formData)
				});
				const data = await response.json();
				if (data.status === 'success') {
					Notiflix.Notify.success(data.message || 'Account created!');
					setTimeout(() => {
						window.location.href = data.redirect || '/';
					}, 800);
				} else {
					Notiflix.Notify.failure(data.message || 'Registration failed');
				}
			} catch (error) {
				Notiflix.Notify.failure('Network error. Please try again.');
			} finally {
				isLoading.value = false;
			}
		};

		return {
			firstName,
			lastName,
			email,
			password,
			passwordConfirm,
			isLoading,
			handleSubmit
		};
	}
}).mount('#app');
