import React, { useState } from 'react';
import { useNavigate,Link } from 'react-router-dom';
import {z} from "zod"

import { useAuth } from '@/context/AuthContext';
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from "@/components/ui/button";

const formSchema = z.object({
    identifier : z.string().min(4,"Enter the username or email"),
    password : z.string().min(6,"password must contain 6 characters").max(30,"password must be at most 30 characters")
});

const LoginPage: React.FC = () => {
    const [identifier, setIdentifier] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [error, setError] = useState<string>('');
    const {login} = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        const result = formSchema.safeParse({ identifier, password });
        if (!result.success) {
            setError('Invalid input');
            return;
        }
        try {
            await login(identifier, password);
            navigate('/home');
        }catch (err) {
            console.error('Error during login:', err);
            setError('Failed to Login, Try again');
        }
    };

    return (
        <div className="flex items-center justify-center min-h-[calc(100vh-64px)]">
            <Card className="w-full max-w-md ">
                <CardHeader>
                    <CardTitle className="text-2xl">Login</CardTitle>
                    <CardDescription>Enter your username or Email and password to log in.</CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className="grid gap-4">
                        <div className="grid gap-2">
                            <Label htmlFor="identifier">Username or Email</Label>
                            <Input id="identifier" type="text" placeholder=" Username or Email" value={identifier} onChange={(e) => setIdentifier(e.target.value)} required/>
                        </div>
                        <div className="grid gap-2">
                            <Label htmlFor="password">Password</Label>
                            <Input id="password" type="password" value={password} placeholder=' Password ' onChange={(e) => setPassword(e.target.value)} required/>
                        </div>
                        {error && <p className="text-red-500 text-sm">{error}</p>}
                        <Button type="submit" className="w-full">
                            Sign in
                        </Button>
                    </form>
                    <div className="mt-4 text-center text-sm">
                        Don't have an account?{' '}
                        <Link to="/signup" className="underline">
                            Sign up
                        </Link>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};

export default LoginPage;
