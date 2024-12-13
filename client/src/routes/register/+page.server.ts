import type { Actions } from './$types'
import { fail } from '@sveltejs/kit'
import { z, ZodError } from "zod";
import { registerUser } from "$lib/api";
import { registerSchema } from "$lib/schemas";
import {handleRegisterError} from "$lib/errors";

export const actions: Actions = {
    default: async ({request}) => {
        const data = await request.formData();
        let formEntries = Object.fromEntries(data)
        try {
            const registerData = registerSchema.parse(formEntries)
            const response = await registerUser(registerData)

            if (response.ok) {
                const json = await response.json();
                console.log(json)
                // TODO store JWTs in cookies


                // TODO call /api/auth/login to auto login the user
                // TODO route to dashboard page
            } else {
                const error = await response.json();
                console.error('Error registering user:', error);
            }

            // return any errors
        } catch (error: unknown) {
            if (error instanceof ZodError) {
                const registrationError = handleRegisterError(error, formEntries)
                console.log(error)
                return fail(400, registrationError.data)
            } else {
                console.error('Error making the request:', error) // TODO: handle error, on +page.svelte
            }
            throw error // rethrow to enable redirect
        }
    }
} satisfies Actions

