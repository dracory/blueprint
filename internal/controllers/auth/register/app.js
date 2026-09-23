const { createApp, ref } = Vue;

createApp({
	setup() {
		// APP_NAME, COUNTRIES, INITIAL_DATA, REGISTER_AJAX_URL and
		// TIMEZONES_AJAX_URL are injected by the server as globals above this
		// script (const, not window properties — reference them directly).
		const appName = typeof APP_NAME !== 'undefined' ? APP_NAME : 'App';
		const countries = typeof COUNTRIES !== 'undefined' ? COUNTRIES : [];
		const initial = typeof INITIAL_DATA !== 'undefined' ? INITIAL_DATA : {};

		const form = ref({
			email: initial.email || '',
			first_name: initial.first_name || '',
			last_name: initial.last_name || '',
			business_name: initial.business_name || '',
			phone: initial.phone || '',
			country: initial.country || '',
			timezone: initial.timezone || '',
		});

		const timezones = ref([]);
		const timezonesLoading = ref(false);
		const isLoading = ref(false);

		const postForm = async (url, fields) => {
			const response = await fetch(url, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/x-www-form-urlencoded',
				},
				body: new URLSearchParams(fields)
			});
			return response.json();
		};

		const loadTimezones = async (country, keepSelected) => {
			const selected = keepSelected ? form.value.timezone : '';
			form.value.timezone = '';
			timezones.value = [];
			if (!country) {
				return;
			}
			timezonesLoading.value = true;
			try {
				const data = await postForm(TIMEZONES_AJAX_URL, { country: country });
				if (data.status === 'success') {
					timezones.value = data.timezones || [];
					if (selected && timezones.value.includes(selected)) {
						form.value.timezone = selected;
					}
				} else {
					Notiflix.Notify.failure(data.message || 'Failed to load timezones');
				}
			} catch (error) {
				Notiflix.Notify.failure('Network error. Please try again.');
			} finally {
				timezonesLoading.value = false;
			}
		};

		const handleCountryChange = () => {
			loadTimezones(form.value.country, false);
		};

		const handleSave = async () => {
			isLoading.value = true;
			try {
				const data = await postForm(REGISTER_AJAX_URL, {
					email: form.value.email,
					first_name: form.value.first_name,
					last_name: form.value.last_name,
					business_name: form.value.business_name,
					phone: form.value.phone,
					country: form.value.country,
					timezone: form.value.timezone,
				});
				if (data.status === 'success') {
					Notiflix.Notify.success('Registration completed!');
					window.location.href = data.redirect || '/user';
				} else {
					Notiflix.Notify.failure(data.message || 'Could not save your details');
				}
			} catch (error) {
				Notiflix.Notify.failure('Network error. Please try again.');
			} finally {
				isLoading.value = false;
			}
		};

		// Preload timezones when the user already has a country selected
		if (form.value.country) {
			loadTimezones(form.value.country, true);
		}

		return {
			appName,
			form,
			countries,
			timezones,
			timezonesLoading,
			isLoading,
			handleCountryChange,
			handleSave,
		};
	}
}).mount('#app');
