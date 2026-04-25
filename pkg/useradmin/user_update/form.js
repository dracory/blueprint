const { createApp } = Vue;

createApp({
    data() {
        return {
            loading: true,
            saving: false,
            action: '',
            errorMessage: '',
            successMessage: '',
            redirectTo: '',
            userId: '',
            returnUrl: '',
            form: {
                status: '',
                first_name: '',
                last_name: '',
                email: '',
                business_name: '',
                phone: '',
                country: '',
                timezone: '',
                memo: ''
            },
            originalEmail: '',
            countries: [],
            timezones: [],
            fieldStatus: {
                first_name: true,
                last_name: true,
                email: true,
                business_name: true,
                phone: true
            }
        };
    },
    computed: {
        hasUnreadableFields() {
            return Object.values(this.fieldStatus).some(v => !v);
        },
        warningFields() {
            const fields = [];
            if (!this.fieldStatus.first_name) fields.push('first name');
            if (!this.fieldStatus.last_name) fields.push('last name');
            if (!this.fieldStatus.email) fields.push('email');
            if (!this.fieldStatus.business_name) fields.push('business name');
            if (!this.fieldStatus.phone) fields.push('phone');
            return fields;
        }
    },
    mounted() {
        this.userId = USER_ID_PLACEHOLDER;
        this.returnUrl = RETURN_URL_PLACEHOLDER;
        this.loadUser();
    },
    methods: {
        async loadUser() {
            this.loading = true;
            this.errorMessage = '';
            try {
                const params = new URLSearchParams({
                    action: 'get-user',
                    user_id: this.userId
                });
                const response = await fetch(urlGetUser + '?' + params.toString());
                const result = await response.json();

                if (result.success) {
                    const d = result.data;
                    this.form.status = d.status || '';
                    this.form.first_name = d.first_name || '';
                    this.form.last_name = d.last_name || '';
                    this.form.email = d.email || '';
                    this.form.business_name = d.business_name || '';
                    this.form.phone = d.phone || '';
                    this.form.country = d.country || '';
                    this.form.timezone = d.timezone || '';
                    this.form.memo = d.memo || '';
                    this.originalEmail = d.email || '';
                    this.countries = d.countries || [];
                    this.timezones = d.timezones || [];
                    if (d.field_status) {
                        this.fieldStatus = { ...this.fieldStatus, ...d.field_status };
                    }
                } else {
                    this.errorMessage = result.message || 'Failed to load user';
                }
            } catch (err) {
                console.error('Error loading user:', err);
                this.errorMessage = 'Failed to load user';
            } finally {
                this.loading = false;
            }
        },
        async onCountryChange() {
            this.form.timezone = '';
            if (!this.form.country) {
                this.timezones = [];
                return;
            }
            try {
                const params = new URLSearchParams({
                    action: 'get-timezones',
                    country_code: this.form.country
                });
                const response = await fetch(urlGetTimezones + '?' + params.toString());
                const result = await response.json();
                if (result.success) {
                    this.timezones = result.data.timezones || [];
                }
            } catch (err) {
                console.error('Error loading timezones:', err);
            }
        },
        async save(actionType) {
            this.action = actionType;
            this.errorMessage = '';
            this.successMessage = '';
            this.redirectTo = '';

            if (!this.form.status) {
                this.errorMessage = 'Status is required';
                return;
            }
            if (!this.form.first_name.trim()) {
                this.errorMessage = 'First name is required';
                return;
            }
            if (!this.form.last_name.trim()) {
                this.errorMessage = 'Last name is required';
                return;
            }
            if (!this.form.email.trim()) {
                this.errorMessage = 'Email is required';
                return;
            }
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(this.form.email.trim())) {
                this.errorMessage = 'Invalid email address';
                return;
            }
            if (!this.form.country) {
                this.errorMessage = 'Country is required';
                return;
            }
            if (!this.form.timezone) {
                this.errorMessage = 'Timezone is required';
                return;
            }

            this.saving = true;
            try {
                const payload = {
                    user_id: this.userId,
                    status: this.form.status.trim(),
                    first_name: this.form.first_name.trim(),
                    last_name: this.form.last_name.trim(),
                    email: this.form.email.trim(),
                    business_name: this.form.business_name.trim(),
                    phone: this.form.phone.trim(),
                    country: this.form.country,
                    timezone: this.form.timezone,
                    memo: this.form.memo.trim()
                };

                const response = await fetch(urlUpdateUser, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });
                const result = await response.json();

                if (result.success) {
                    this.successMessage = 'User saved successfully';
                    if (actionType === 'save') {
                        this.redirectTo = this.returnUrl;
                        Swal.fire({
                            title: 'Saved!',
                            text: 'User has been saved successfully.',
                            icon: 'success',
                            timer: 3000,
                            timerProgressBar: true,
                            showConfirmButton: false
                        });
                        setTimeout(() => {
                            if (this.redirectTo) {
                                window.location.href = this.redirectTo;
                            }
                        }, 3000);
                    } else {
                        Swal.fire({
                            title: 'Saved!',
                            text: 'User has been saved successfully.',
                            icon: 'success',
                            toast: true,
                            position: 'top-end',
                            timer: 3000,
                            timerProgressBar: true,
                            showConfirmButton: false
                        });
                    }
                } else {
                    this.errorMessage = result.message || 'Failed to save user';
                }
            } catch (err) {
                console.error('Error saving user:', err);
                this.errorMessage = 'Failed to save user';
            } finally {
                this.saving = false;
            }
        }
    }
}).mount('#app-user-update');
