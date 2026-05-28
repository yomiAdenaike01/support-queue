import { Outlet } from "react-router-dom";
import { Header } from "@/components/layout/Header";
import { Sidebar } from "@/components/layout/Sidebar";

export function Layout() {
  return (
    <div className="min-h-screen bg-[--color-bg]">
      <Sidebar />
      <div className="pl-20 lg:pl-64">
        <Header />
        <main className="p-4 lg:p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
