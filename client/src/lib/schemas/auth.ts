import {z} from "zod";

export const passwordSchema = z
    .string()
    .min(8, { message: "Password must be at least 8 characters long" })
    .max(20, { message: "Password too long" })
    .refine((password) => /[A-Z]/.test(password), {
        message: "Password must contain at least 1 uppercase",
    })
    .refine((password) => /[a-z]/.test(password), {
        message: "Password must contain at least 1 lowercase",
    })
    .refine((password) => /[0-9]/.test(password),
        { message: "Password must contain a number" })
    .refine((password) => /[!@#$%^&*]/.test(password), {
        message: "Password must contain a special character (!@#$%^&*)",
    });

export const registerSchema = z.object({
    username: z.string().min(3, 'username must be at least 3 characters long'),
    email: z.string().email('Invalid email address'),
    password: passwordSchema,
    confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
    message: "Passwords must match",
    path: ["confirmPassword"],
})