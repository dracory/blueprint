const { createApp, ref } = Vue;

createApp({
	setup() {
		const email = ref('');
		const step = ref(1);
		const isLoading = ref(false);

		const handleSendLink = async () => {
			if (!email.value) {
				Notiflix.Notify.failure('Please enter your email');
				return;
			}

			isLoading.value = true;

			const formData = new FormData();
			formData.append('email', email.value);
			if (RETURN_URL) {
				formData.append('return', decodeURIComponent(RETURN_URL));
			}

			try {
				const response = await fetch(SEND_AJAX_URL, {
					method: 'POST',
					headers: {
						'Content-Type': 'application/x-www-form-urlencoded',
					},
					body: new URLSearchParams(formData)
				});
				const data = await response.json();
				if (data.status === 'success') {
					Notiflix.Notify.success('Link sent! Check your email.');
					step.value = 2;
				} else {
					Notiflix.Notify.failure(data.message || 'Failed to send link');
				}
			} catch (error) {
				Notiflix.Notify.failure('Network error. Please try again.');
			} finally {
				isLoading.value = false;
			}
		};

		const handleBack = () => {
			step.value = 1;
		};

		return {
			email,
			step,
			isLoading,
			handleSendLink,
			handleBack
		};
	}
}).mount('#app');
