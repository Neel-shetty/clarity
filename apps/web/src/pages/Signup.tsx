import React , { useState } from "react";
import { useNavigate,Link } from "react-router-dom";
import axios from "axios";
import { z } from "zod";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";

const formSchema = z.object({
    username: z.string().min(5,"Username must be atleast 5 characters long.").regex(/^[a-zA-Z0-9_-]+$/, "Username can only contain letters, numbers, underscores, and hyphens."),
    email : z.string().email("Invalid email"),
    password: z.string().min(6, "Password must be at least 6 characters long."),
    confirmPassword: z.string().min(6,"please confirm your password"),
});

export default function Signup() {
    const navigate = useNavigate();
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [email,setEmail] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [error, setError] = useState("");

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (password !== confirmPassword) {
            setError('Passwords do not match.');
            return;
        }
        const result = formSchema.safeParse({ username,email,password, confirmPassword });
        if (!result.success) {
            const firstError = result.error.issues[0]?.message || "Invalid input.";
            setError(firstError);
            return;
        }
        try {
            const response = await axios.post(
                "/user/signup",
                { username,email,password },
                {
                    headers: { "Content-Type": "application/json" },
                    withCredentials: true
                }
            );
            if (response.status === 200) {
                navigate("/");
            } else {
                setError("Failed to sign up. Please try again.");
            }
        } catch (err) {
            console.error("Error signing up:", err);
            setError("Failed to sign up. Please try again.");
        }
    };
    return <div className="flex items-center justify-center min-h-[calc(100vh-64px)]">
        <Card className="w-full max-w-md ">
            <CardHeader className="text-center">
                <CardTitle className="text-2xl">Signup</CardTitle>
                <CardDescription>Enter your username, Email and password to signup.</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4">
                <form className="grid gap-4" onSubmit={handleSubmit}> 
                    <div className="grid gap-2">
                        <Label htmlFor="username">Username</Label>
                        <Input id="username" type="text" placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} required />
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="email">Email</Label>
                        <Input id="email" type="email" placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} required />
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="password">Password</Label>
                        <Input id="password" type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} required/>
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="confirmPassword">Confirm Password</Label>
                        <Input id="confirmPassword" type="password" value={confirmPassword} onChange={e=> setConfirmPassword(e.target.value)} placeholder="Confirm Password" required/>
                    </div>
                    <Button type="submit">Signup</Button>
                    <div className="mt-4 text-center text-sm">
                        Already have an account?{' '}
                    </div>
                    <Link to="/" className="underline text-center ">
                        Login
                    </Link>
                    {error && <p className="text-red-500 text-center">{error}</p>}
                </form>
            </CardContent>
        </Card>
    </div>
}