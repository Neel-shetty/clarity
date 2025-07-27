import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/button";

export default function Home() {
    const {user,logout} = useAuth();
    return (
        <div>
            <div className="flex items-center justify-between px-6 py-4 bg-gray-900 text-white shadow">
                <span className="text-lg font-semibold">
                    Hello, {user.name}
                </span>
                <Button variant="outline" onClick={logout}>
                    Logout
                </Button>
            </div>
            <div className="p-6">
                <h1>Hello World</h1>
            </div>
        </div>
    );
}