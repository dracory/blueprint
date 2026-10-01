const { createApp, ref } = Vue;

createApp({
	setup() {
		const email = ref('');
		const password = ref('');
		const isLoading = ref(false);

		const handleLogin = async () => {
			if (!email.value || !password.value) {
				Notiflix.Notify.failure('Please enter your email and password');
				return;
			}

			isLoading.value = true;

			const formData = new FormData();
			formData.append('email', email.value);
			formData.append('password', password.value);
			if (RETURN_URL) {
				formData.append('return', decodeURIComponent(RETURN_URL));
			}

			try {
				const response = await fetch(LOGIN_AJAX_URL, {
					method: 'POST',
					headers: {
						'Content-Type': 'application/x-www-form-urlencoded',
					},
					body: new URLSearchParams(formData)
				});
				const data = await response.json();
				if (data.status === 'success') {
					Notiflix.Notify.success('Login successful!');
					window.location.href = data.redirect || '/';
				} else {
					Notiflix.Notify.failure(data.message || 'Login failed');
				}
			} catch (error) {
				Notiflix.Notify.failure('Network error. Please try again.');
			} finally {
				isLoading.value = false;
			}
		};

		return {
			email,
			password,
			isLoading,
			handleLogin
		};
	}
}).mount('#app');
