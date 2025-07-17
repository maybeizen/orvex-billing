import Hero from "@/components/hero";
import Pricing from "@/components/pricing";
import About from "@/components/about";
import FAQ from "@/components/faq";
import Contact from "@/components/contact";
import Services from "@/components/services";

export default function Home() {
  return (
    <div>
      <Hero />
      <Pricing />
      <Services />
      <About />
      <FAQ />
      <Contact />
    </div>
  );
}
