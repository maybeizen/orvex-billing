import React from "react";
import LightRays from "../components/ui/anim/light-rays";
import About from "../components/about";
import Pricing from "../components/pricing";
import FAQ from "../components/faq";
import Footer from "../components/footer";
import { Button } from "../components/ui/button";
import TextType from "../components/ui/anim/text-type";

export default function Home() {
  const features = [
    {
      title: "Lightning Fast",
      description:
        "Deploy your applications in seconds with our optimized build pipeline and global CDN distribution.",
      icon: "far fa-bolt",
    },
    {
      title: "Reliable",
      description:
        "99.9% uptime guarantee with automatic failover and monitoring to keep your apps running smoothly.",
      icon: "far fa-check-circle",
    },
    {
      title: "Developer Friendly",
      description:
        "Intuitive APIs, comprehensive documentation, and seamless integration with your existing workflow.",
      icon: "far fa-code",
    },
    {
      title: "Scalable",
      description:
        "Auto-scaling infrastructure that grows with your application from prototype to enterprise scale.",
      icon: "far fa-chart-line",
    },
    {
      title: "Analytics",
      description:
        "Real-time insights into your application performance with detailed metrics and user analytics.",
      icon: "far fa-chart-bar",
    },
    {
      title: "Secure",
      description:
        "Enterprise-grade security with SSL certificates, DDoS protection, and secure environment variables.",
      icon: "fas fa-shield-alt",
    },
  ];
  return (
    <div className="min-h-screen bg-neutral-950 text-white relative">
      <div className="relative z-10 flex items-center justify-center min-h-screen">
        <div className="absolute inset-0 z-0">
          <LightRays
            raysOrigin="top-center"
            raysColor="#ffffff"
            raysSpeed={0.4}
            lightSpread={1}
            rayLength={1.2}
            followMouse={false}
            mouseInfluence={0.1}
            noiseAmount={0.1}
            distortion={0.05}
            className="w-full h-full"
          />
        </div>
        <div className="text-center flex flex-col items-center justify-center font-semibold z-10">
          <h1 className="text-5xl text-white">Minecraft Hosting</h1>
          <div className="flex flex-row items-center justify-center text-5xl">
            <div>
              <span className="text-white">Made</span>{" "}
              <span className="text-violet-500">
                <TextType
                  text={["Better", "Simpler", "More Affordable"]}
                  textColors={["#8e51ff"]}
                  typingSpeed={75}
                  pauseDuration={1500}
                  showCursor={true}
                  cursorCharacter="_"
                  className="text-violet-500 mb-4"
                />
              </span>
            </div>
          </div>
          {/* <SplitText
            text="Made Better"
            className="text-5xl font-light text-center mb-4 text-violet-500"
            delay={100}
            duration={0.6}
            ease="power3.out"
            splitType="chars"
            from={{ opacity: 0, y: 40 }}
            to={{ opacity: 1, y: 0 }}
            threshold={0.1}
            rootMargin="-100px"
            textAlign="center"
            onLetterAnimationComplete={handleAnimationComplete}
          /> */}
          <p className="text-xl font-normal text-neutral-400 w-xl mb-8">
            Experience next-level Minecraft hosting from just $2.80. Launch your
            server in minutes—no hassle, just play.
          </p>
          <div className="flex flex-row justify-center items-center gap-4">
            <Button variant="white" rounded="lg" size="md">
              Start Playing
            </Button>
            <Button variant="glass" rounded="lg" size="md">
              Learn More
            </Button>
          </div>
        </div>
      </div>

      <div className="relative z-10 py-20 px-8">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-semibold mb-4 text-transparent bg-clip-text bg-gradient-to-r from-neutral-400 via-white to-neutral-300">
              Powerful Features
            </h2>
            <p className="text-lg font-normal text-neutral-400 max-w-2xl mx-auto">
              Everything you need to streamline your development process and
              ship faster.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((feature, index) => (
              <div
                key={index}
                className="p-6 bg-neutral-900/50 backdrop-blur-sm border border-white/10 rounded-3xl"
              >
                <div className="flex flex-col h-full">
                  <div className="w-10 h-10 bg-white/10 rounded-lg flex items-center justify-center mb-4">
                    <i
                      className={`${feature.icon} text-violet-400 text-lg`}
                    ></i>
                  </div>
                  <h3 className="text-lg font-normal mb-2 text-white">
                    {feature.title}
                  </h3>
                  <p className="text-neutral-400 text-sm font-normal leading-relaxed">
                    {feature.description}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      <About />
      <Pricing />
      <FAQ />
      <Footer />
    </div>
  );
}
