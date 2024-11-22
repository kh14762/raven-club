<script>
    import { onMount } from 'svelte';

    let isDarkMode = false;

    onMount(() => {
        // Check if user has a saved preference
        const savedTheme = localStorage.getItem('theme');

        // Check if user has system dark mode preference
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

        // Set initial theme based on saved preference or system preference
        isDarkMode = savedTheme ? savedTheme === 'dark' : prefersDark;

        // Apply theme immediately on mount
        applyTheme(isDarkMode);
    });

    function toggleTheme() {
        isDarkMode = !isDarkMode;
        applyTheme(isDarkMode);
        localStorage.setItem('theme', isDarkMode ? 'dark' : 'light');
    }

    function applyTheme(dark) {
        // Ensure we're updating the document class
        if (typeof window !== 'undefined') {
            if (dark) {
                document.documentElement.classList.add('dark');
            } else {
                document.documentElement.classList.remove('dark');
            }
        }
    }
</script>

<button
        class="relative inline-flex items-center h-6 w-11 rounded-full bg-gray-200 dark:bg-gray-700 transition-colors duration-300"
        on:click={toggleTheme}
        aria-label="Toggle theme"
>
  <span
          class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform duration-300"
          class:translate-x-6={isDarkMode}
          class:translate-x-1={!isDarkMode}
  >
    <!-- Sun icon for light mode -->
    <svg
            class="h-4 w-4 text-yellow-500"
            class:hidden={isDarkMode}
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
    >
      <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707"
      />
    </svg>
      <!-- Moon icon for dark mode -->
    <svg
            class="h-4 w-4 text-gray-400"
            class:hidden={!isDarkMode}
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
    >
      <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
      />
    </svg>
  </span>
</button>