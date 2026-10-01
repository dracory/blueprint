const { createApp, ref } = Vue;

createApp({
	setup() {
		const password = ref('');
		const passwordConfirm = ref('');
		const step = ref(1);
		const isLoading = ref(false);
		const loginUrl = LOGIN_URL;

		const handleSubmit = async () => {
			if (!password.value || !passwordConfirm.value) {
				Notiflix.Notify.failure('Please enter and confirm your new password');
				return;
			}

			if (password.value !== passwordConfirm.value) {
				Notiflix.Notify.failure('Passwords do not match');
				return;
			}

			isLoading.value = true;

			const formData = new FormData();
			formData.append('token', TOKEN);
			formData.append('password', password.value);
			formData.append('password_confirm', passwordConfirm.value);

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
					Notiflix.Notify.success(data.message || 'Password reset successfully!');
					step.value = 2;
					setTimeout(() => {
						window.location.href = data.redirect || loginUrl;
					}, 1500);
				} else {
					Notiflix.Notify.failure(data.message || 'Failed to reset password');
				}
			} catch (error) {
				Notiflix.Notify.failure('Network error. Please try again.');
			} finally {
				isLoading.value = false;
			}
		};

		return {
			password,
			passwordConfirm,
			step,
			isLoading,
			loginUrl,
			handleSubmit
		};
	}
}).mount('#app');
