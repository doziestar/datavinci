import React from "react";
import { motion, Variants } from "framer-motion";

const HoverAnimated3DHeading: React.FC = () => {
  const headingVariants: Variants = {
    initial: { opacity: 1, y: 0 },
    hover: {
      scale: 1.05,
      rotateX: 10,
      rotateY: 10,
      transition: { duration: 0.3, ease: "easeInOut" },
    },
  };

  const letterVariants: Variants = {
    initial: { y: 0 },
    hover: {
      y: [-5, 5, -5],
      transition: { repeat: Infinity, duration: 1, ease: "easeInOut" },
    },
  };

  const text = "Empower Your Data with DataVinci";

  return (
    <motion.h1
      className="lg:leading-tighter text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl xl:text-7xl bg-clip-text text-transparent bg-gradient-to-r from-primary to-secondary perspective-1000"
      variants={headingVariants}
      initial="initial"
      whileHover="hover"
      style={{ transformStyle: "preserve-3d" }}
    >
      {text.split(" ").map((word, wordIndex) => (
        <span
          key={wordIndex}
          style={{ display: "inline-block", whiteSpace: "nowrap" }}
        >
          {word.split("").map((char, charIndex) => (
            <motion.span
              key={`${wordIndex}-${charIndex}`}
              variants={letterVariants}
              style={{ display: "inline-block" }}
              custom={(wordIndex * word.length + charIndex) * 0.1}
              transition={{
                delay: (wordIndex * word.length + charIndex) * 0.05,
              }}
            >
              {char}
            </motion.span>
          ))}
          {wordIndex < text.split(" ").length - 1 && "\u00A0"}
        </span>
      ))}
    </motion.h1>
  );
};

export default HoverAnimated3DHeading;
