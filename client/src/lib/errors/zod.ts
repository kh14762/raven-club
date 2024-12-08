import { ZodError } from "zod";
import {fail} from "@sveltejs/kit";

export function handleRegisterError(error: ZodError, formEntries: Record<string, FormDataEntryValue>) {
    const errors = error.flatten()
    const {username, email, password, confirmPassword} = formEntries
    const {fieldErrors} = errors
    return fail(400, {
        username: username,
        email: email,
        password: password,
        confirmPassword: confirmPassword,
        error: {
            ...(fieldErrors?.username ? {field: 'username', message: fieldErrors.username[0]} : {}),
            ...(fieldErrors?.email ? {field: 'email', message: fieldErrors.email[0]} : {}),
            ...(fieldErrors?.password ? {field: 'password', message: fieldErrors.password[0]} : {}),
            ...(fieldErrors?.confirmPassword ? {
                field: 'confirm',
                message: fieldErrors.confirmPassword[0]
            } : {})
        }
    })
}