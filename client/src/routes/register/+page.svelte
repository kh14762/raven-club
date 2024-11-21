<script lang="ts">
    let username = '';
    let email = '';
    let password = '';
    let confirmPassword = '';
    let isLoading = false;
    let errorMessage = '';

    async function handleRegister(event: SubmitEvent) {
        event.preventDefault();
        errorMessage = '';

        // Basic validation
        if (password !== confirmPassword) {
            errorMessage = 'Passwords do not match';
            return;
        }

        isLoading = true;
        try {
            // Add your registration logic here
            console.log('Register attempt:', { username, email, password });

            // After successful registration, you might want to redirect
            // await goto('/login');
        } catch (error) {
            console.error('Registration failed:', error);
            errorMessage = 'Registration failed. Please try again.';
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow-md">
        <div>
            <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
                Create your account
            </h2>
        </div>

        {#if errorMessage}
            <div class="bg-red-50 text-red-500 p-3 rounded-md text-sm">
                {errorMessage}
            </div>
        {/if}

        <form class="mt-8 space-y-6" on:submit={handleRegister}>
            <div class="rounded-md shadow-sm space-y-4">
                <!-- Username field -->
                <div>
                    <label for="username" class="sr-only">Username</label>
                    <input
                            id="username"
                            name="username"
                            type="text"
                            bind:value={username}
                            required
                            class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                            placeholder="Username"
                    />
                </div>

                <!-- Email field -->
                <div>
                    <label for="email" class="sr-only">Email address</label>
                    <input
                            id="email"
                            name="email"
                            type="email"
                            bind:value={email}
                            required
                            class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                            placeholder="Email address"
                    />
                </div>

                <!-- Password field -->
                <div>
                    <label for="password" class="sr-only">Password</label>
                    <input
                            id="password"
                            name="password"
                            type="password"
                            bind:value={password}
                            required
                            class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                            placeholder="Password"
                    />
                </div>

                <!-- Confirm Password field -->
                <div>
                    <label for="confirm-password" class="sr-only">Confirm Password</label>
                    <input
                            id="confirm-password"
                            name="confirm-password"
                            type="password"
                            bind:value={confirmPassword}
                            required
                            class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                            placeholder="Confirm Password"
                    />
                </div>
            </div>

            <div>
                <button
                        type="submit"
                        disabled={isLoading}
                        class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    {#if isLoading}
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              Loading...
            </span>
                    {/if}
                    Create Account
                </button>
            </div>

            <div class="text-center">
                <p class="mt-2 text-sm text-gray-600">
                    Already have an account?
                    <a href="/login" class="font-medium text-indigo-600 hover:text-indigo-500">
                        Sign in
                    </a>
                </p>
            </div>
        </form>
    </div>
</div>