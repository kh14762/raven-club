<script lang="ts">
    let email = '';
    let password = '';
    let isLoading = false;
    let errorMessage = '';

    async function handleSubmit() {
        errorMessage = '';
        isLoading = true;

        try {
            // Add your login logic here
            console.log('Login attempt:', { email, password });

            // After successful login, you might want to redirect
            // await goto('/dashboard');
        } catch (error) {
            console.error('Login failed:', error);
            errorMessage = 'Invalid email or password';
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow-md">
        <div>
            <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
                Sign in to your account
            </h2>
        </div>

        {#if errorMessage}
            <div class="bg-red-50 text-red-500 p-3 rounded-md text-sm">
                {errorMessage}
            </div>
        {/if}

        <form class="mt-8 space-y-6" on:submit|preventDefault={handleSubmit}>
            <div class="rounded-md shadow-sm space-y-4">
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
            </div>

            <div class="flex items-center justify-between">
                <div class="flex items-center">
                    <input
                            id="remember-me"
                            name="remember-me"
                            type="checkbox"
                            class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                    />
                    <label for="remember-me" class="ml-2 block text-sm text-gray-900">
                        Remember me
                    </label>
                </div>

                <div class="text-sm">
                    <a href="/forgot-password" class="font-medium text-indigo-600 hover:text-indigo-500">
                        Forgot your password?
                    </a>
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
                    Sign in
                </button>
            </div>

            <div class="text-center">
                <p class="mt-2 text-sm text-gray-600">
                    Don't have an account?
                    <a href="/register" class="font-medium text-indigo-600 hover:text-indigo-500">
                        Sign up
                    </a>
                </p>
            </div>
        </form>
    </div>
</div>