<script lang="ts">
    import { onMount } from 'svelte';
    import { tweened } from 'svelte/motion';
    import { linear } from 'svelte/easing';

    import darkStars from '$lib/icons/dark-stars-thirteen.png'
    import lightStars from '$lib/icons/light-stars-thirteen.png'
    import darkRavenIcon from '$lib/icons/raven-dark.png';
    import lightRavenIcon from '$lib/icons/raven-light.png'

    let rotationSpeed = 3000; // Default speed (3 seconds per rotation)
    let isRotating = true;
    let currentRotation = 0;

    // Create the tweened store with initial configuration
    const rotation = tweened(0, {
        duration: rotationSpeed,
        easing: linear
    });

    function rotate() {
        if (!isRotating) return;

        currentRotation += 360;
        // Update the duration each time we set a new rotation
        rotation.set(currentRotation, {
            duration: rotationSpeed,
            easing: linear
        }).then(rotate);
    }

    // Function to update speed
    function updateSpeed(newSpeed: number) {
        rotationSpeed = newSpeed;

        // Store current progress
        const currentValue = $rotation;

        // Stop current animation
        isRotating = false;

        // Update tweened store with new duration
        rotation.set(currentValue, { duration: 0 }).then(() => {
            // Restart animation with new speed
            isRotating = true;
            currentRotation = currentValue;
            rotate();
        });
    }

    function toggleRotation() {
        isRotating = !isRotating;
        if (isRotating) {
            currentRotation = $rotation;
            rotate();
        }
    }

    onMount(() => {
        rotate();
        return () => {
            isRotating = false;
        };
    });
</script>

<div class="min-h-screen flex flex-col items-center justify-center bg-gray-50 dark:bg-gray-900 transition-colors duration-300">
    <div class="relative">
        <!-- Stationary base image -->
        <img
                src={darkRavenIcon}
                alt="Stationary element"
                class="w-full h-full object-contain rotate-12"
        />
            <!-- Rotating overlay image -->
        <img
                src={lightStars}
                alt="Rotating element"
                class="absolute top-0 left-0 w-full h-full object-contain animate-spin aspect-square"
                style="animation: spin 15s linear infinite;"
        />
    </div>
</div>

<style>

</style>