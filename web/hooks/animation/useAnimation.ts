import { Variants, useSpring } from "framer-motion";

export const useAnimations = () => {
  const jelly: Variants = {
    hover: {
      scale: [1, 1.25, 0.75, 1.15, 0.95, 1.05, 1],
      rotate: [0, 5, -5, 3, -3, 1, 0],
      transition: { duration: 0.6 },
    },
    tap: {
      scale: [1, 1.05, 0.95, 1],
      boxShadow: [
        "0 0 0 0 rgba(0, 0, 0, 0)",
        "0 0 0 10px rgba(0, 0, 0, 0.1)",
        "0 0 0 20px rgba(0, 0, 0, 0.1)",
        "0 0 0 0 rgba(0, 0, 0, 0)",
      ],
      transition: { duration: 0.4 },
    },
  };

  const splash: Variants = {
    tap: {
      scale: [1, 1.05, 0.95, 1],
      boxShadow: [
        "0 0 0 0 rgba(0, 0, 0, 0)",
        "0 0 0 10px rgba(0, 0, 0, 0.1)",
        "0 0 0 20px rgba(0, 0, 0, 0.1)",
        "0 0 0 0 rgba(0, 0, 0, 0)",
      ],
      transition: { duration: 0.4 },
    },
  };

  const smoothAppear: Variants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.5, ease: "easeOut" },
    },
  };

  const wobble = useSpring(0, {
    stiffness: 300,
    damping: 10,
  });

  return { jelly, splash, smoothAppear, wobble };
};
