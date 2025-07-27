import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';

export default function NotFound() {
    return (
            <div className="flex flex-col items-center justify-center min-h-screen">
                <h1 className="text-6xl font-bold text-white-800 mb-4">404</h1>
                <p className="text-xl text-white-600 mb-8">Oops! Page Not Found.</p>
                <Button asChild>
                    <Link to="/">Go to Home</Link>
                </Button>
            </div>
    );
}