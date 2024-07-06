import React, { useEffect } from "react";
import { motion, useAnimationControls } from "framer-motion";

const Dancing3DHeading = () => {
  const controls = useAnimationControls();

  const headingVariants = {
    hidden: { opacity: 0, y: 50 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.8, ease: "easeOut" },
    },
    hover: {
      scale: 1.05,
      rotateX: 10,
      rotateY: 10,
      transition: { duration: 0.3, ease: "easeInOut" },
    },
  };

  useEffect(() => {
    controls.start((i) => ({
      y: [0, -10, 0],
      transition: {
        repeat: Infinity,
        repeatType: "mirror",
        duration: 1 + i * 0.1,
        ease: "easeInOut",
      },
    }));
  }, [controls]);

  const text = "Empower Your Data with DataVinci";

  return (
    <motion.h1
      className="lg:leading-tighter text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl xl:text-7xl bg-clip-text text-transparent bg-gradient-to-r from-primary to-secondary perspective-1000"
      variants={headingVariants}
      initial="hidden"
      animate="visible"
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
              animate={controls}
              custom={(wordIndex * word.length + charIndex) * 0.1}
              style={{ display: "inline-block" }}
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

export default Dancing3DHeading;
