mod call;
pub mod context;
mod genesis;

#[cfg(test)]
mod tests;

#[cfg(feature = "native")]
mod query;

pub use call::CallMessage;
use context::TransferContext;
#[cfg(feature = "native")]
pub use query::Response;
use sov_modules_api::{Error, ModuleInfo};
use sov_state::WorkingSet;

pub struct TransferConfig {}

#[cfg_attr(feature = "native", derive(sov_modules_api::ModuleCallJsonSchema))]
#[derive(ModuleInfo, Clone)]
pub struct Transfer<C: sov_modules_api::Context> {
    /// Address of the module.
    #[address]
    pub address: C::Address,

    /// Reference to the Bank module.
    #[module]
    pub(crate) bank: sov_bank::Bank<C>,

    /// Keeps track of the address of each token we minted.
    /// The index is the hash of the token denom (using the hasher `C::Hasher`).
    /// Note: we use `Vec<u8>` instead of `Output<C::Hasher>` because `C::Hasher`
    /// is not cloneable, and we currently need our module to be cloneable
    #[state]
    pub(crate) minted_tokens: sov_state::StateMap<Vec<u8>, C::Address>,

    /// Keeps track of the address of each token we escrowed as a function of
    /// the (hash of) the token denom. We need this map because we have the
    /// token address information when escrowing the tokens (i.e. when someone
    /// calls a `send_transfer()`), but not when unescrowing tokens (i.e in a
    /// `recv_packet`), in which case the only information we have is the ICS 20
    /// denom, and amount. Given that every token that is unescrowed has been
    /// previously escrowed, our strategy to get the token address associated
    /// with a denom is
    /// 1. when tokens are escrowed, save the mapping `denom -> token address`
    /// 2. when tokens are unescrowed, lookup the token address by `denom`
    ///
    /// Note: Even though we could store the `denom: String` as a key, we prefer
    /// to hash it to the key a constant size.
    #[state]
    pub(crate) escrowed_tokens: sov_state::StateMap<Vec<u8>, C::Address>,
}

impl<C: sov_modules_api::Context> Transfer<C> {
    pub fn into_context<'ws, 'c>(
        self,
        sdk_context: &'c C,
        working_set: &'ws mut WorkingSet<C::Storage>,
    ) -> TransferContext<'ws, 'c, C> {
        TransferContext::new(self, sdk_context, working_set)
    }
}

impl<C: sov_modules_api::Context> sov_modules_api::Module for Transfer<C> {
    type Context = C;

    type Config = TransferConfig;

    type CallMessage = call::CallMessage;

    fn genesis(
        &self,
        config: &Self::Config,
        working_set: &mut WorkingSet<C::Storage>,
    ) -> Result<(), Error> {
        // The initialization logic
        Ok(self.init_module(config, working_set)?)
    }

    fn call(
        &self,
        _msg: Self::CallMessage,
        _context: &Self::Context,
        _working_set: &mut WorkingSet<C::Storage>,
    ) -> Result<sov_modules_api::CallResponse, Error> {
        todo!()
    }
}

impl<C> core::fmt::Debug for Transfer<C>
where
    C: sov_modules_api::Context,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        // FIXME: put real values here, or remove `Debug` requirement from router::Module
        f.debug_struct("Transfer")
            .field("address", &self.address)
            .finish()
    }
}
