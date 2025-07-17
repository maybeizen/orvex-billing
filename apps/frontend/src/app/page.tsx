import Hero from "@/components/hero";
import Pricing from "@/components/pricing";
import About from "@/components/about";
import FAQ from "@/components/faq";
import Contact from "@/components/contact";
import Services from "@/components/services";
import Navbar from "@/components/navbar";
import Footer from "@/components/footer";

export default function Home() {
  return (
    <div>
      <Navbar />
      <Hero />
      <Pricing />
      <Services />
      <About />
      <FAQ />
      <Contact />
      <Footer />
    </div>
  );
}
