import { useCallback, useState } from "react";
import {
  Button,
  TextField,
  Paper,
  Typography,
  Box,
  Stack,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import * as Yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup/src/yup.js";

const signupInitialValue = {
  username: "",
  email: "",
  password: "",
  userRole: "seller",
};

const loginInitialValue = {
  email: "",
  password: "",
  userRole: "seller",
};

function generateSchema(isSignUp: boolean) {
  return Yup.object().shape({
    ...(isSignUp && {
      username: Yup.string().required("Username is required"),
    }),
    email: Yup.string().email("Invalid email").required("Email is required"),
    password: Yup.string()
      .required("Password is required")
      .matches(
        /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/,
        "Password must contain at least 8 characters, one uppercase letter, one lowercase letter, one number, and one special character"
      ),
    userRole: Yup.string().required("User role is required"),
  });
}

export const Login = () => {
  const navigation = useNavigate();
  const [isSignUp, setIsSignUp] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  const {
    register,
    handleSubmit,
    watch,
    reset,
    formState: { errors },
  } = useForm<any>({
    defaultValues: isSignUp ? signupInitialValue : loginInitialValue,
    mode: "all",
    resolver: yupResolver(generateSchema(isSignUp)),
  });

  const handleLogin = useCallback(async () => {
    try {
      const response = await fetch("http://localhost:8080/v1/log-in", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: watch("email"),
          password: watch("password"),
          role: watch("userRole"),
        }),
      });
      console.log({ response: response.body });
      if (response.status === 200) {
        const data = await response.json();
        localStorage.setItem("token", data.token);
        localStorage.setItem("userRole", watch("userRole"))
        reset();
        navigation("/dashboard");
        setError("");
      }else{
        setError("Please enter correct email password");
      }
    } catch (error) {
      console.log(error);
     
    }
  }, [watch, reset, navigation]);

  const handleSignUp = useCallback(async () => {
    try {
      const response = await fetch("http://localhost:8080/v1/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: watch("email"),
          password: watch("password"),
          name: watch("username"),
          phoneNo: watch("phone"),
          role: watch("userRole"),
        }),
      });

      if (response.status === 200) {
        setIsSignUp(false);
        reset();
        setError("")
      }else{
        setError("Something went wrong")
      }
    } catch (error) {
      console.log(error);
    }
  }, [watch, reset]);

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100vh",
      }}
    >
      <Paper elevation={3} sx={{ padding: 4, width: 300 }}>
        <Typography variant="h5" mb={2} textAlign="center">
          {!isSignUp ? "Login" : "Sign-Up"}
        </Typography>
        <Stack
          component={"form"}
          spacing={2}
          onSubmit={
            isSignUp ? handleSubmit(handleSignUp) : handleSubmit(handleLogin)
          }
        >
          {isSignUp && (
            <TextField
              label="Name"
              variant="outlined"
              fullWidth
              helperText={errors?.username?.message as string}
              {...register("username")}
              sx={{ "& .MuiFormHelperText-root": { color: "red" } }}
            />
          )}
          <TextField
            label="Email"
            variant="outlined"
            fullWidth
            helperText={errors?.email?.message as string}
            {...register("email")}
            sx={{ "& .MuiFormHelperText-root": { color: "red" } }}
          />
          <TextField
            label="Password"
            type="password"
            variant="outlined"
            fullWidth
            helperText={errors?.password?.message as string}
            {...register("password")}
            sx={{ "& .MuiFormHelperText-root": { color: "red" } }}
          />
          <FormControl fullWidth>
            <InputLabel id="userRole">User Role</InputLabel>
            <Select
              labelId="userRole"
              id="userRole"
              label="User Role"
              value={watch("userRole")}
              {...register("userRole")}
              sx={{ "& .MuiFormHelperText-root": { color: "red" } }}
            >
              <MenuItem value={"buyer"}>Buyer</MenuItem>
              <MenuItem value={"seller"}>Seller</MenuItem>
            </Select>
          </FormControl>
          {error && <Typography color="error">{error}</Typography>}
          <Button variant="contained" color="primary" type="submit" fullWidth>
            {!isSignUp ? "Login" : "SignUp"}
          </Button>
        </Stack>
        <Stack direction="row" justifyContent="space-between">
          <Button
            variant="text"
            sx={{
              textAlign: "flex-end",
            }}
            onClick={() => {
              setIsSignUp((prev) => !prev);
              reset();
              setError("");
            }}
          >
            {isSignUp ? "Login" : "SignUp"}
          </Button>
        </Stack>
      </Paper>
    </Box>
  );
};
