const { createApp, ref } = Vue;

createApp({
	setup() {
		const email = ref('');
		const firstName = ref('');
		const step = ref(1);
		const isLoading = ref(false);
		const loginUrl = LOGIN_URL;

		const handleSubmit = async () => {
			if (!email.value || !firstName.value) {
				Notiflix.Notify.failure('Please enter your email and first name');
				return;
			}

			isLoading.value = true;

			const formData = new FormData();
			formData.append('email', email.value);
			formData.append('first_name', firstName.value);

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
					Notiflix.Notify.success(data.message || 'If the account details match, a reset link has been sent.');
					step.value = 2;
				} else {
					Notiflix.Notify.failure(data.message || 'Failed to process request');
				}
			} catch (error) {
				Notiflix.Notify.failure('Network error. Please try again.');
			} finally {
				isLoading.value = false;
			}
		};

		return {
			email,
			firstName,
			step,
			isLoading,
			loginUrl,
			handleSubmit
		};
	}
}).mount('#app');
