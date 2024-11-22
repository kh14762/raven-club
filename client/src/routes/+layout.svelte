<script lang="ts">
    import '../styles/app.css'
    import ThemeToggle from "$lib/ThemeToggle.svelte";

    let isMenuOpen = false;

    function toggleMenu() {
        isMenuOpen = !isMenuOpen;
    }
</script>

<header class="bg-white dark:bg-gray-800 shadow-sm dark:shadow-gray-700 transition-colors duration-300">
    <nav class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
            <!-- Left side - Logo/Brand -->
            <div class="flex items-center">
                <a href="/" class="text-xl font-bold text-indigo-600">
                    Raven Club
                </a>
            </div>

            <!-- Desktop Navigation -->
            <div class="hidden sm:flex sm:items-center sm:space-x-4">
                <a
                        href="/login"
                        class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 dark:text-gray-200
                               hover:text-gray-900 dark:hover:text-white hover:bg-gray-50 dark:hover:bg-gray-700"
                >
                    Login
                </a>
                <a
                        href="/register"
                        class="px-3 py-2 rounded-md text-sm font-medium bg-indigo-600 text-white
                               hover:bg-indigo-700 dark:bg-indigo-500 dark:hover:bg-indigo-600"
                >
                    Register
                </a>
                <ThemeToggle/>
            </div>

            <!-- Mobile menu button -->
            <div class="flex items-center sm:hidden">
                <button
                        type="button"
                        aria-controls="mobile-menu"
                        class="inline-flex items-center justify-center p-2 rounded-md
                               text-gray-400 dark:text-gray-300
                               hover:text-gray-500 dark:hover:text-gray-100
                               hover:bg-gray-100 dark:hover:bg-gray-700"
                        aria-expanded="false"
                        on:click={toggleMenu}
                >
                    <span class="sr-only">Open main menu</span>
                    <!-- Icon when menu is closed -->
                    {#if !isMenuOpen}
                        <svg
                                class="block h-6 w-6"
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                                aria-hidden="true"
                        >
                            <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M4 6h16M4 12h16M4 18h16"
                            />
                        </svg>
                    {:else}
                        <!-- Icon when menu is open -->
                        <svg
                                class="block h-6 w-6"
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                                aria-hidden="true"
                        >
                            <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    {/if}
                </button>
            </div>
        </div>

        <!-- Mobile menu -->
        {#if isMenuOpen}
            <div class="sm:hidden" id="mobile-menu">
                <div class="pt-2 pb-3 space-y-1">
                    <a
                            href="/login"
                            class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 dark:text-gray-200
                                   hover:text-gray-900 dark:hover:text-white hover:bg-gray-50 dark:hover:bg-gray-700"
                    >
                        Login
                    </a>
                    <a
                            href="/register"
                            class="btn-primary"
                    >
                        Register
                    </a>
                    <ThemeToggle/>
                </div>
            </div>
        {/if}
    </nav>
</header>
<slot />
