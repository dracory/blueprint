const { createApp, ref, computed, nextTick } = Vue;

createApp({
	setup() {
		const email = ref('');
		const otpDigits = ref(['', '', '', '', '', '']);
		const step = ref(1);
		const isLoading = ref(false);
		const nonce = ref('');

		const otp = computed(() => otpDigits.value.join(''));

		const handleOtpInput = (event, index) => {
			const value = event.target.value;
			// Only allow digits
			if (value && !/^\d$/.test(value)) {
				otpDigits.value[index] = '';
				event.target.value = '';
				return;
			}
			if (value && index < 5) {
				const inputs = document.querySelectorAll('.otp-input');
				if (inputs[index + 1]) {
					inputs[index + 1].focus();
				}
			}
		};

		const handleOtpPaste = (event, index) => {
			event.preventDefault();
			const pasteData = (event.clipboardData || window.clipboardData).getData('text');
			const digits = pasteData.replace(/\D/g, '').split('');

			if (digits.length === 0) return;

			// Fill from current index onwards
			for (let i = 0; i < digits.length && (index + i) < 6; i++) {
				otpDigits.value[index + i] = digits[i];
			}

			// Focus the appropriate input after paste
			nextTick(() => {
				const inputs = document.querySelectorAll('.otp-input');
				const focusIndex = Math.min(index + digits.length, 5);
				if (inputs[focusIndex]) {
					inputs[focusIndex].focus();
				}
			});
		};

		const handleOtpDelete = (event, index) => {
			if (!otpDigits.value[index] && index > 0) {
				const inputs = document.querySelectorAll('.otp-input');
				if (inputs[index - 1]) {
					inputs[index - 1].focus();
				}
			}
		};

		const handleSendOtp = async () => {
			if (!email.value) {
				Notiflix.Notify.failure('Please enter your email');
				return;
			}

			isLoading.value = true;

			const formData = new FormData();
			formData.append('email', email.value);

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
					Notiflix.Notify.success('Code sent! Check your email.');
					nonce.value = data.nonce || '';
					step.value = 2;
					nextTick(() => {
						const inputs = document.querySelectorAll('.otp-input');
						if (inputs[0]) inputs[0].focus();
					});
				} else {
					Notiflix.Notify.failure(data.message || 'Failed to send code');
				}
			} catch (error) {
				Notiflix.Notify.failure('Network error. Please try again.');
			} finally {
				isLoading.value = false;
			}
		};

		const handleVerifyOtp = async () => {
			if (otp.value.length !== 6) {
				Notiflix.Notify.failure('Please enter the 6-digit code');
				return;
			}

			isLoading.value = true;

			const formData = new FormData();
			formData.append('email', email.value);
			formData.append('otp', otp.value);
			formData.append('nonce', nonce.value);

			try {
				const response = await fetch(VERIFY_AJAX_URL, {
					method: 'POST',
					headers: {
						'Content-Type': 'application/x-www-form-urlencoded',
					},
					body: new URLSearchParams(formData)
				});
				const data = await response.json();
				if (data.status === 'success') {
					window.location.href = data.redirect || '/user';
				} else {
					Notiflix.Notify.failure(data.message || 'Invalid code');
				}
			} catch (error) {
				Notiflix.Notify.failure('Network error. Please try again.');
			} finally {
				isLoading.value = false;
			}
		};

		const handleBack = () => {
			step.value = 1;
			otpDigits.value = ['', '', '', '', '', ''];
			nonce.value = '';
		};

		return {
			email,
			otpDigits,
			otp,
			nonce,
			step,
			isLoading,
			handleSendOtp,
			handleVerifyOtp,
			handleBack,
			handleOtpInput,
			handleOtpPaste,
			handleOtpDelete
		};
	}
}).mount('#app');
